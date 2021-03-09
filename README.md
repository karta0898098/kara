# kara
kara 是封裝了一些第三方 package 的 utils，讓建立新服務時簡單調用 Constructor 即可。</br>

### db package
提供了 gorm client，並統一使用 [zerolog]( https://github.com/rs/zerolog "link") 作為 gorm 的 logger，以及初始化時的 retry 機制防止 db 忙碌不能連線而導致 server panic 的情況。

### errors package
使用了 [pkg/errors](https://github.com/pkg/errors "link") 來封裝 error 並預設提供了一些錯誤。讓 err 可在 middleware 中做錯誤處理即可。

### http package
提供以下 feature
* [dump log](https://github.com/karta0898098/kara/blob/master/http/middleware/dump.go "link") 
* [error record](https://github.com/karta0898098/kara/blob/master/http/middleware/error.go "link")
* [http log](https://github.com/karta0898098/kara/blob/master/http/middleware/logger.go "link")
* [trace id](https://github.com/karta0898098/kara/blob/master/http/middleware/tracer.go "link")


### zlog package
提供初始化 log 的方法，並預設提供了設定，log level , format 等相關設定。


### grpc package
 // TODO