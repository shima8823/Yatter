# Yatter
Twitter・Mastodonライクな仮想のサービスのBackend API

 - Go言語を用いたbackendの実装経験を得たい
 - ソフトウェアテストの実践経験を得たい
 - ソフトウェアのアーキテクチャについて学びたい

【頑張ったところ/工夫したところ】
この課題では初期の環境が提供されており、APIとそのテストの開発を行いました。主にアカウントの操作（登録、取得）、投稿のCRUD、タイムラインの取得、フォロー関連の操作、そしてそのテストを実装しました。APIの動作確認には、Swagger UIを使用して詳細なテストを行いました。

特にテスト部分に注力しました。Webアプリのテスト実装は初めての経験でした。初期段階でDatabase Access Object（以下dao）のテストにモックを使用しましたが、http handlerのテストでは同じモックを使うとテストが重複する問題があったため、DAOでは実DB、handlerではモックを使用する形にしました。
また、HTTPテストの冗長性を解消するため、サブテストからテーブル駆動テストへ移行し、コードの可読性を向上させる工夫をしました。

アーキテクチャ設計にも注力しました。handlerがdaoに直接依存せず、repositoryインターフェースを介してアクセスするようにしたことで、daoとhandlerを疎結合に保ちました。これにより、テストが書きやすくなり、コードも変更しやすい構造を持つようになりました。このアプローチは依存性逆転の原則に基づいており、その原則を学ぶ機会ともなりました。

このプロジェクトを通じて、Goの文法、テスト方法、そしてアーキテクチャの基礎を学びました。改善点として挙げられるのは、テストがモックを多用してしまった結果、コードロジックとテストロジックの区別が曖昧になってしまった点です。今後はDBの操作を直接モックするのではなく、dao層そのものをモック化することで、コードの可読性を向上させると考えています。
