## 1 cors
[HTTP 访问控制](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/CORS)

[Official CORS gin's middleware](https://github.com/gin-contrib/cors)

## 2 简单请求
简单请求不会触发 cors 预检请求。满足下述条件的视为简单请求：
- 使用下列方法之一：
    - GET
    - HEAD
    - POST
- 除了被用户代理自动设置的标头字段（例如 Connection、User-Agent 或其他在 Fetch 规范中定义为禁用标头名称的标头），允许人为设置的字段为 Fetch 规范定义的对 CORS 安全的标头字段集合。该集合为： 
    - Accept
    - Accept-Language
    - Content-Language
    - Content-Type（需要注意额外的限制）
    - Range（只允许简单的范围标头值 如 bytes=256- 或 bytes=127-255）
- Content-Type 标头所指定的媒体类型的值仅限于下列三者之一： 
    - text/plain
    - multipart/form-data
    - application/x-www-form-urlencoded
    - 如果请求是使用 XMLHttpRequest 对象发出的，在返回的 XMLHttpRequest.upload 对象属性上没有注册任何事件监听器；也就是说，给定一个 XMLHttpRequest 实例 xhr，没有调用 xhr.upload.addEventListener()，以监听该上传请求。
- 请求中没有使用 ReadableStream 对象。

**示例**：
假如站点 `https://foo.example` 的网页应用想要访问 `https://bar.other` 的资源。foo.example 的网页中可能包含类似于下面的 JavaScript 代码：
```js
const xhr = new XMLHttpRequest();
const url = 'https://bar.other/resources/public-data/';

xhr.open('GET', url);
xhr.onreadystatechange = someHandler;
xhr.send();
```

请求报文如下：
```http
GET /resources/public-data/ HTTP/1.1
Host: bar.other
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:71.0) Gecko/20100101 Firefox/71.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
Accept-Language: en-us,en;q=0.5
Accept-Encoding: gzip,deflate
Connection: keep-alive
Origin: https://foo.example
```

请求标头字段 `Origin` 表明该请求来源于 `http://foo.example`。

响应报文如下：
```http
HTTP/1.1 200 OK
Date: Mon, 01 Dec 2008 00:23:53 GMT
Server: Apache/2
Access-Control-Allow-Origin: *
Keep-Alive: timeout=2, max=100
Connection: Keep-Alive
Transfer-Encoding: chunked
Content-Type: application/xml

[…XML Data…]
```

本例中，服务端返回的 `Access-Control-Allow-Origin` 标头的 `Access-Control-Allow-Origin: *` 值表明，该资源可以被任意外源访问。
```http
Access-Control-Allow-Origin: *
```

使用 `Origin` 和 `Access-Control-Allow-Origin` 就能完成最简单的访问控制。如果 `https://bar.other` 的资源持有者想限制他的资源只能通过 `https://foo.example` 来访问（也就是说，非 `https://foo.example` 域无法通过跨源访问访问到该资源），他可以这样做：
```http
Access-Control-Allow-Origin: https://foo.example
```

> 注意：当响应的是附带身份凭证的请求时，服务端必须明确 Access-Control-Allow-Origin 的值，而不能使用通配符“*”。

## 3 预检请求
与简单请求不同，“需预检的请求”要求必须首先使用 OPTIONS 方法发起一个预检请求到服务器，以获知服务器是否允许该实际请求。“预检请求”的使用，可以避免跨域请求对服务器的用户数据产生未预期的影响。

如下是一个需要执行预检请求的 HTTP 请求：
```js
const xhr = new XMLHttpRequest();
xhr.open('POST', 'https://bar.other/resources/post-here/');
xhr.setRequestHeader('X-PINGOTHER', 'pingpong');
xhr.setRequestHeader('Content-Type', 'application/xml');
xhr.onreadystatechange = handler;
xhr.send('<person><name>Arun</name></person>');
```

上面的代码使用 `POST` 请求发送一个 XML 请求体，该请求包含了一个非标准的 HTTP `X-PINGOTHER` 请求标头。这样的请求标头并不是 HTTP/1.1 的一部分，但通常对于 web 应用很有用处。另外，该请求的 `Content-Type` 为 `application/xml`，且使用了自定义的请求标头，所以该请求需要首先发起“预检请求”。

下面是服务端和客户端完整的信息交互。首次交互是预检请求/响应：
```http
OPTIONS /doc HTTP/1.1
Host: bar.other
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:71.0) Gecko/20100101 Firefox/71.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
Accept-Language: en-us,en;q=0.5
Accept-Encoding: gzip,deflate
Connection: keep-alive
Origin: https://foo.example
Access-Control-Request-Method: POST
Access-Control-Request-Headers: X-PINGOTHER, Content-Type

HTTP/1.1 204 No Content
Date: Mon, 01 Dec 2008 01:15:39 GMT
Server: Apache/2
Access-Control-Allow-Origin: https://foo.example
Access-Control-Allow-Methods: POST, GET, OPTIONS
Access-Control-Allow-Headers: X-PINGOTHER, Content-Type
Access-Control-Max-Age: 86400
Vary: Accept-Encoding, Origin
Keep-Alive: timeout=2, max=100
Connection: Keep-Alive
```

从上面的报文中，我们看到，第 1 - 10 行使用 OPTIONS 方法发送了预检请求，浏览器根据上面的 JavaScript 代码片断所使用的请求参数来决定是否需要发送，这样服务器就可以回应是否可以接受用实际的请求参数来发送请求。OPTIONS 是 HTTP/1.1 协议中定义的方法，用于从服务器获取更多信息，是安全的方法。该方法不会对服务器资源产生影响。注意 OPTIONS 预检请求中同时携带了下面两个标头字段：
```http
Access-Control-Request-Method: POST
Access-Control-Request-Headers: X-PINGOTHER, Content-Type
```

标头字段 `Access-Control-Request-Method` 告知服务器，实际请求将使用 POST 方法。标头字段 `Access-Control-Request-Headers` 告知服务器，实际请求将携带两个自定义请求标头字段：`X-PINGOTHER` 与 `Content-Type`。服务器据此决定，该实际请求是否被允许。

第 12 - 21 行为预检请求的响应，表明服务器将接受后续的实际请求方法（POST）和请求头（X-PINGOTHER）。重点看第 15 - 18 行：
```http
Access-Control-Allow-Origin: https://foo.example
Access-Control-Allow-Methods: POST, GET, OPTIONS
Access-Control-Allow-Headers: X-PINGOTHER, Content-Type
Access-Control-Max-Age: 86400
```

服务器的响应携带了 `Access-Control-Allow-Origin: https://foo.example`，从而限制请求的源域。同时，携带的 `Access-Control-Allow-Methods` 表明服务器允许客户端使用 `POST` 和 `GET` 方法发起请求（与 `Allow` 响应标头类似，但该标头具有严格的访问控制）。

标头字段 `Access-Control-Allow-Headers` 表明服务器允许请求中携带字段 `X-PINGOTHER` 与 `Content-Type`。与 `Access-Control-Allow-Methods` 一样，`Access-Control-Allow-Headers` 的值为逗号分割的列表。

最后，标头字段 `Access-Control-Max-Age` 给定了该预检请求可供缓存的时间长短，单位为秒，默认值是 5 秒。在有效时间内，浏览器无须为同一请求再次发起预检请求。以上例子中，该响应的有效时间为 86400 秒，也就是 24 小时。请注意，浏览器自身维护了一个最大有效时间，如果该标头字段的值超过了最大有效时间，将不会生效。

预检请求完成之后，发送实际请求：
```http
POST /doc HTTP/1.1
Host: bar.other
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:71.0) Gecko/20100101 Firefox/71.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
Accept-Language: en-us,en;q=0.5
Accept-Encoding: gzip,deflate
Connection: keep-alive
X-PINGOTHER: pingpong
Content-Type: text/xml; charset=UTF-8
Referer: https://foo.example/examples/preflightInvocation.html
Content-Length: 55
Origin: https://foo.example
Pragma: no-cache
Cache-Control: no-cache

<person><name>Arun</name></person>

HTTP/1.1 200 OK
Date: Mon, 01 Dec 2008 01:15:40 GMT
Server: Apache/2
Access-Control-Allow-Origin: https://foo.example
Vary: Accept-Encoding, Origin
Content-Encoding: gzip
Content-Length: 235
Keep-Alive: timeout=2, max=99
Connection: Keep-Alive
Content-Type: text/plain

[Some XML payload]
```

## 4 HTTP 请求标头字段
本节列出了可用于发起跨源请求的标头字段。请注意，这些标头字段无须手动设置。当开发者使用 `XMLHttpRequest` 对象发起跨源请求时，它们已经被设置就绪。

### 4.1 Origin
`Origin` 标头字段表明预检请求或实际跨源请求的源站。
```http
Origin: <origin>
```

origin 参数的值为源站 URL。它不包含任何路径信息，只是服务器名称（可以为null）。**注意，在所有访问控制请求中，Origin 标头字段总是被发送。**

### 4.2 Access-Control-Request-Method
`Access-Control-Request-Method` 标头字段用于预检请求。其作用是，将实际请求所使用的 HTTP 方法告诉服务器。
```http
Access-Control-Request-Method: <method>
```

### 4.3 Access-Control-Request-Headers
`Access-Control-Request-Headers` 标头字段用于预检请求。其作用是，将实际请求所携带的标头字段（通过 `setRequestHeader()` 等设置的）告诉服务器。这个浏览器端标头将由互补的服务器端标头 `Access-Control-Allow-Headers` 回答。
```http
Access-Control-Request-Headers: <field-name>[, <field-name>]*
```

## 5 HTTP 响应标头字段

### 5.1 Access-Control-Allow-Origin
响应标头中可以携带一个 Access-Control-Allow-Origin 字段，其语法如下：
```http
Access-Control-Allow-Origin: <origin> | *
```

`Access-Control-Allow-Origin` 参数指定了单一的源，告诉浏览器允许该源访问资源。或者，对于不需要携带身份凭证的请求，服务器可以指定该字段的值为通配符“`*`”，表示允许来自任意源的请求。

例如，为了允许来自 https://mozilla.org 的代码访问资源，你可以指定：
```http
Access-Control-Allow-Origin: https://mozilla.org
Vary: Origin
```

如果服务端指定了具体的单个源（作为允许列表的一部分，可能会根据请求的来源而动态改变）而非通配符“`*`”，那么响应标头中的 `Vary` 字段的值必须包含 `Origin`。这将告诉客户端：服务器对不同的 `Origin` 返回不同的内容。

### 5.2 Access-Control-Expose-Headers
译者注：在跨源访问时，`XMLHttpRequest` 对象的 `getResponseHeader()` 方法只能拿到一些最基本的响应头，Cache-Control、Content-Language、Content-Type、Expires、Last-Modified、Pragma，如果要访问其他头，则需要服务器设置本响应头。

例如：
```http
Access-Control-Expose-Headers: X-My-Custom-Header, X-Another-Custom-Header
```

这样浏览器就能够通过 `getResponseHeader` 访问 `X-My-Custom-Header` 和 `X-Another-Custom-Header` 响应头了。

### 5.3 Access-Control-Max-Age
`Access-Control-Max-Age` 头指定了 `preflight` 请求的结果能够被缓存多久。
```http
Access-Control-Max-Age: <delta-seconds>
```

`delta-seconds` 参数表示 preflight 预检请求的结果在多少秒内有效。

### 5.4 Access-Control-Allow-Credentials
`Access-Control-Allow-Credentials` 头指定了当浏览器的 `credentials` 设置为 true 时是否允许浏览器读取 response 的内容。当用在对 preflight 预检测请求的响应中时，它指定了实际的请求是否可以使用 `credentials`。请注意：简单 `GET` 请求不会被预检；如果对此类请求的响应中不包含该字段，这个响应将被忽略掉，并且浏览器也不会将相应内容返回给网页。
```http
Access-Control-Allow-Credentials: true
```

### 5.5 Access-Control-Allow-Methods
`Access-Control-Allow-Methods` 标头字段指定了访问资源时允许使用的请求方法，用于预检请求的响应。其指明了实际请求所允许使用的 HTTP 方法。
```http
Access-Control-Allow-Methods: <method>[, <method>]*
```

### 5.6 Access-Control-Allow-Headers
`Access-Control-Allow-Headers` 标头字段用于预检请求的响应。其指明了实际请求中允许携带的标头字段。这个标头是服务器端对浏览器端 `Access-Control-Request-Headers` 标头的响应。
```http
Access-Control-Allow-Headers: <header-name>[, <header-name>]*
```

