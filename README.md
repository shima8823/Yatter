# Yatter

### Set Up
```sh
docker compose up
make
make test
```

## function
#### アカウント
 - POST /v1/accounts<br>
 - GET /v1/accounts/username<br>

#### 投稿
 - GET /v1/statuses/id<br>
 - DELETE /statuses/id<br>
 - パブリックタイムラインの取得<br>
GET /v1/timelines/public<br>

#### フォロー関連機能
 - POST /accounts/username/follow<br>
 - GET /accounts/username/following<br>
 - GET /accounts/username/followers<br>
 - アカウントのunfollow<br>
POST /accounts/username/unfollow<br>
 - アカウントとのrelation取得<br>
GET/accounts/relationships<br>
 - home timeline取得<br>
GET /timelines/home<br>
