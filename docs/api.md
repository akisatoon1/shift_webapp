# api仕様

## ベースURL
`http://localhost:3000/api`

## Header
- `Cookie: <cookie-key>=<cookie-value>`
- `Content-Type: application/json`

## 認証,認可
ほぼ全てのエンドポイント(GET /loginを除く)でCookieが必要.

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
**ログインを試みる**
- 成功時: `200 OK`
- 失敗時: `401 Unauthorized`
#### Request Body
```
{
    "login_id": string,
    "password": string
}
```

### GET /session
**セッション情報を返す**
**ログインしていないときは、`401 Unauthorized`エラーを返す**
#### Response Body
```
{
    "user": {
        "id": number,
        "name": string,
        "roles": string[],
        "created_at": string
    }
}
```


### DELETE /session
**ログアウトする**
#### Response
`200 OK`

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

### GET /requests/{request_id}
**提出されたシフトエントリーの一覧をを含む、シフトリクエスト詳細データを返す**
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
    "created_at": string,

    "submissions": {
        "submitter": {
            "id": number,
            "name": string
        }
    }[],

    "entries": {
        "id": number,
        "user": {
            "id": number,
            "name": string
        },
        "date": string,
        "hour": number
    }[]
}
```

### GET /requests/{request_id}/submissions/mine
**自分のシフト提出を取得**
#### Response body
```
{
    "submission": {
        "id": number,

        "user": {
            "id": number,
            "name": string
        },

        "entries": [
            {
                "id": number,
                "date": string,
                "hour": number
            }
        ]
    } | null
}
```

### POST /requests/{request_id}/submissions
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
    "id": number // 提出id
}
```
