# Message API

## 概要 (Overview)

本プロジェクトのメッセージ機能のバックエンド処理を担当する、Go言語で実装されたAPIです。メインのLaravelアプリケーションから呼び出されることを想定しています。

## APIエンドポイント (API Endpoints)

本APIが提供するエンドポイントは以下の通りです。

| Method | Endpoint | 説明 |
| :--- | :--- | :--- |
| `POST` | `/messages/register` | APIを利用するクライアント（アプリケーションユーザー）を新規登録します。 |
| `POST` | `/messages/token` | リフレッシュトークンを使い、アクセストークンを再発行します。 |
| `GET` | `/messages/` | メッセージの一覧を取得します。 |
| `POST` | `/messages/` | 新しいメッセージを作成します。 |
| `PATCH` | `/messages/{id}` | 指定したメッセージの内容を更新します。 |
| `DELETE` | `/messages/{id}` | 指定したメッセージを削除します。 |
| `POST` | `/messages/send-scheduled` | （メッセージ予約送信トリガーシステムからの呼び出し用）送信予約されたメッセージのステータスを「送信予定」から「送信済み」に更新します。|

---

## 📗 Supplementary Material

-   [**設計に関する意思決定の記録 (Design Decisions)**](./DESIGN_DECISIONS.md)