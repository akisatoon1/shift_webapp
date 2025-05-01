# api仕様

## ベースURL
`http://localhost:3000/api`

## Header
- `Authorization`
- `Content-Type: application/json`

## 認証
ほぼ全てのエンドポイント(GET /loginを除く)でトークンが必要.

`Authorization: Session {token}`

## エラー
エラーメッセージを返す.
```
{
    "error": string
}
```

## ステータスコード
### 成功時
- GET: `200`
- POST: `201`

## エンドポイント一覧

### POST /login
**ログイン成功したら、Session Tokenを返す**
#### Request Body
```
{
    "login_id": string,
    "password": string
}
```
#### Response Body
```
{
    "session_token": string,
    "expires_at": number
}
```

### DELETE /session
**ログアウトする**
#### Response
`204 No Content`

### GET /requests
**リクエスト一覧を返す**
#### Response body
```
{
    "id": number,
    "creator": {
        "id": number,
        "name": string
    },
    "start_date": string,
    "end_date": string,
    "deadline": string,
    "created_at": string
}[]
```

### POST /requests
**新しいリクエストを追加し、新しいIDを返す**
#### Request body
```
{
    "start_date": string  // 開始日
    "end_date": string    // 終了日
    "deadline": string    // 提出の期限
}
```
#### Reponse body
```
{
    "id": number
}
```

### GET /requests/{request_id}/entries
**提出されたシフトエントリーの一覧を返す**
#### Response body
```
{
    "id": number,
    "entries": {
        "id": number,
        "user": {
            "id": number,
            "name": string
        },
        "date": string,
        "hour": number
    }[]
}[]
```

### POST /requests/{request_id}/entries
**新しいシフトエントリーを提出(追加)して、新しいIDを返す**
#### Request body
```
{
    "date": string,  // シフトに入る日付
    "hour": number   // シフトに入る時刻
}[]
```
#### Response body
```
{
    "id": number,
    "entries": {
        "id": number
    }[]
}
```
