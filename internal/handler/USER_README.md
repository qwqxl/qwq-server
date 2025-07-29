### ✅ 登录请求
发送用户名、密码和设备类型：

```js
// 假设用 fetch，或 axios 类似
const login = async () => {
  const res = await fetch("http://localhost:8080/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      username: "alice",
      password: "123456",
      device: "pc", // 或 mobile
    }),
  });

  const data = await res.json();

  // 保存 token 到本地存储
  localStorage.setItem("access_token", data.access_token);
  localStorage.setItem("refresh_token", data.refresh_token);
};

```


### ✅ 请求受保护接口（自动刷新）
每次请求前端接口时，要把 Token 带上：

```js
const apiRequest = async (url: string, method = "GET", body?: any) => {
  const headers: any = {
    "Content-Type": "application/json",
    "Authorization": localStorage.getItem("access_token") || "",
    "X-Refresh-Token": localStorage.getItem("refresh_token") || "",
    "X-Device": "pc", // 与登录时保持一致
  };

  const res = await fetch(`http://localhost:8080${url}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  // 如果服务器返回了新的 token（刷新后的）
  const newAccess = res.headers.get("Authorization");
  const newRefresh = res.headers.get("X-Refresh-Token");
  if (newAccess && newRefresh) {
    localStorage.setItem("access_token", newAccess);
    localStorage.setItem("refresh_token", newRefresh);
  }

  if (res.status === 401) {
    alert("登录过期，请重新登录");
    // 可能跳转登录页
  }

  return await res.json();
};

```

#### 调用示例

```js
const getProfile = async () => {
  const data = await apiRequest("/profile");
  console.log(data);
};

```

### ✅ 登出接口

```js
const logout = async () => {
  const res = await fetch("http://localhost:8080/logout", {
    method: "POST",
    headers: {
      "Authorization": localStorage.getItem("access_token") || "",
      "X-Refresh-Token": localStorage.getItem("refresh_token") || "",
      "X-Device": "pc",
    },
  });
  localStorage.removeItem("access_token");
  localStorage.removeItem("refresh_token");
  alert("已退出");
};

```

### ✅ 总结调用流程图

```
[登录页面]
   ↓
POST /login
   → 保存 access_token + refresh_token
   ↓
[调用受保护接口]
   → 自动带上 access_token / refresh_token / X-Device
   ↓
[后端发现 access_token 过期]
   → 返回新的 token → 覆盖保存
   ↓
[继续请求成功，用户无感知]

```