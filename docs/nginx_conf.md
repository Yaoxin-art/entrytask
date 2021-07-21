# nginx config for profile access

## upstream for web app
```text
    upstream entrytask {
        server 127.0.0.1:7777;
    }
```

## locations for 
```text
        location /user {
            proxy_pass  http://entrytask;
        }

        location / {
            root        /Users/zhenrong.zeng/Workspaces/Test/golang/entrytask/htmls;
            index       login.html;
        }
```
