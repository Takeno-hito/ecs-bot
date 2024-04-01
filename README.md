# ECS Bot

VALORANT Campus Summit のために作成した Discord bot、の外部公開用 repo です。


## お断り

外部公開用に少し内容を削っているので、動かなかったらごめんなさい。  
他の人が使うことを全然考えてなくてごめんなさい。ソースコードを読んで使ってください  
そしてその上でとても人に見せられるソースコードをしてないのでご了承ください…  

## 使い方

1. run.go の `Run()` で使うメソッドを選ぶ
2. あらかじめ用意しておいた選手/チームスプレッドシートなどを csv に出力する
3. bot/data.csv に貼り付ける (data.csv.sample 参照)
4. bot/run.go のカラム名/変数名などをいい感じに編集する （ごめん 頑張ってソースコード読んでください）
5. 実行するとチームだったら チーム Role と チーム vc / text チャンネルが自動生成されます
6. 選手の方だったら選手に自動でロールが振られます
7. 今ならなんと dice bot と (Riot API の申請をもらうための) 二言語対応機能がついてくる！
8. Ctrl+C なりで Bot を落としてください (main.go:32-34 をコメントアウトすれば Run だけ走って落ちます)