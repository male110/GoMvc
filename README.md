<p>
  <b>当前版本：1.08</b>
  <br/>
  <b>最后更新日期：2014-07-25</b>
</p>
<a href="http://www.66a.cm/" target="_blank">官网</a><br/>
<a href="https://github.com/male110/GoMvc/archive/master.zip">下载GoMvc</a><br/>

<p>有任何问题，可加QQ群：184572648，我基本上每天都在线的</p>
<p>原域名因忘记续费被别人抢注了，新域名为66a.cm</p>
<a href="#updatelog">更新日志</a><br/>
<a href="#build"> 编译</a><br />
<a href="#config">  配置文件</a><br />
<a href="#route">  路由注册</a><br />
<a href="#yilai">对其它包的依赖</a>
<p>
  <b>
    <a name="updatelog">更新日志</a>
  </b>
  <pre>
    <b>2014-07-25</b>	    修复RenderAction模板函数Cookies传递的BUG
    <b>2014-06-13</b>	    修改日志记录System/Log/Logger.go,AddError自动记录堆栈信息，增加AddErrMsg函数，自动记录堆栈信息，Add不记堆栈信息。
	2014-06-04	    Controller增加ClearSession函数，RenderView增加错误日志
    2014-05-23	    修改Session相关处理部分，在配置文件中，<SessionType>0</SessionType>配成零或空，表示程序不使用Session，
		    比如做WEBAPI时，可以配置成0，程序不使用Session可以降低资源占用，提高性能。
    2014-05-22 　　 修改Http请求处理过程，支持这样的Action
    　　　　　　　　func (this *Controller) IsExist() string {
    　　　　　　　　}
    　　　　　　　　action可以直接返回一个string类型，如果是其它的非ViewResult类型，会转换成string并输出。
    2014-05-20　　  修改Controller的IsPost属性,修改Binder的错误。
    2014-05-12　　  Controller加入Redirect函数，修改RederAction,修改RederAction。。
    2014-05-08　　　加入模板函数RenderView，更新文档；将System\TemplateFunc包跟System\ViewEngine合并成一个包。<br>
    2014-05-07　　　修改RanderAction模板函数的一处错误。<br>
    2014-05-05　　　修改RanderAction模板函数的Bug，所有错误日志记录堆栈信息，以便调试，处理错误。<br>
        　　　　　　　　增加编译的批处理文件Windows下运行build.bat,Linux下运行build.sh<br/>
    2014-01-24　　　程序意外退出时，记录错误日志。<br/>
    2013-10-14　　　增加域的功能。<br/>
    2013-10-12　　　修复路由和FieSession的Bug。
      </pre>
</p>
<p>
  <a name="build">
    <b>编译</b>
  </a>
  <div>
    GoMVC是一个简单，便捷的MVC框架。程序注释全部使用中文，很适合国人使用。文档也很详细。
    编译时，需要把GoMvc目录设置为GOPATH.
  </div>
</p>
<p>
  <b>
    <a name="config">配置文件</a>
  </b>
</p>
<div>
  <p>
    网站的配置文件为web.config，格式为XML，配置项的内容如下：
  </p>
  <p>
    <b>ShowErrors：</b>是否显示错误信息。true,显示；false,不显示。建义在测试时可以设置为true,发布到正式环境后设置为false。
  </p>
  <p>
    <b>CookieDomain：</b>Cookies的Domain信息，可用来共享cookie。如domain.com，和sub.domain.com，可以通过把CookieDomain统一设置为domain.com来共享cookies信息
  </p>
  <p>
    <b>Theme：</b>网站当前使用的主题，在Views目录下，可以有多套网站模板。
  </p>
  <p>
    <b>LogPath：</b>日志文件的存放位置
  </p>
  <p>
    <b>LogFileMaxSize：</b>单个日志文件的大小，超过指定大小后将创建一个新的日志文件。
  </p>
  <p>
    <b>DriverName：</b>数据库的驱动名称。
  </p>
  <p>
    <b>DataSourceName：</b>数据库的连接字符串。
  </p>
  <p>
    <b>StaticDir：</b>静态目录,该目录下通常是CSS,JS,图片等静态资源。
  </p>
  <p>
    <b>StaticFile：</b>静态文件，用来设置单个的静态文件，主要是为了提高灵活性，满足特殊的需求.
  </p>
  <p>
    <b>SessionType：</b>Session的存放类型,1,文件,2内存,3Mysql数据库,修改需重启才能生效。当配置为3时，需要在数据库中创建一个表，来存放session,创建表的SQL如下：<br />
  </p>
  <pre>
    CREATE TABLE `session` (
    `session_id` CHAR(32) NULL,
    `session_data` BLOB NULL,
    `lastupdatetime` DATETIME NULL,
    PRIMARY KEY (`session_id`)
    )
    COLLATE=&#39;utf8_general_ci&#39;;
  </pre>
  <p>
    <b>SessionLocation：</b>当SessionType为1时，该项为Session文件的存放路径；SessionType为3时,该项为数据库连接字符串。
  </p>
  <p>
    <b>SessionTimeOut：</b>Session超时时间，单位分钟
  </p>
  <p>
    <b>MemFreeInterval：</b>程序中有定时器，定时对Session进行检查，删除超时的Session，该配置项用来设置多久进行一次检查，单位秒，默认值60。
  </p>
  <p>
    <b>ListenPort：</b>网站的端口号,该配置改后必须重启程序才能生效。
  </p>
  <p>
    &nbsp;
  </p>
</div>
<p>
  <b>
    <a name="route">  路由注册</a>
  </b>
</p>
<p>
  用RouteTable.AddRote来注册路由。其格式如下：
</p>
<pre>
  //注册标准路由
  RouteTable.AddRote(&amp;RouteItem{
  Name:     &quot;default&quot;,
  Url:      &quot;{controller}/{action}&quot;,
  Defaults: map[string]interface{}{&quot;controller&quot;: &quot;home&quot;, &quot;action&quot;: &quot;index&quot;}})
</pre>
<p>
  Name:路由名称<br />
  Url:路由的格式<br />
  Defaults: 路由参数的默认值
</p>
除了默认值，还可以指定约束，来限制参数的类型，如下面的例子，指定id参数，只能是数字型。
<pre>
  RouteTable.AddRote(&amp;RouteItem{
  Name:        &quot;default&quot;,
  Url:         &quot;{controller}/{action}/{id}&quot;,
  Defaults:    map[string]interface{}{&quot;controller&quot;: &quot;home&quot;, &quot;action&quot;: &quot;index&quot;, &quot;id&quot;: 123},
  Constraints: map[string]string{&quot;id&quot;: `^(\d+)$`}})
</pre>
在上面的例子中我们指定了id参数只能是数字，并设置了默认值123。要在Controller中获取该参数值，可以用this.RouteData[&quot;id&quot;]。
<p>
  因为在Go没有办法反射出包中的所有struct，所以需要手动来注册Controller,格式如下：
</p>
<pre>
  import (
  &quot;System/Web&quot;
  &quot;fmt&quot;
  )

  type Home struct {
  Web.Controller
  }

  //注册Controller
  func init() {
  Web.App.RegisterController(Home{})
  }
</pre>
对于Controller的命名没有严格的要求，可以用Home,也可以用HomeController


<p>
  <b>
    <a name="yilai">对其它包的依赖</a>
  </b>
</p>
<p>
  GOMvc追求简单，实用，尽可能减少对其它包的依赖。在GoMvc中有两个地方用到了第三方包：
</p>
<p>
  1，System/Session/MysqlSession.go 该文件实现了以mysql的方式来存储Session的功能，需要mysql的驱动。可以此下载<a
        href="https://github.com/go-sql-driver/mysql" target="_blank">https://github.com/go-sql-driver/mysql</a>
</p>
<p>
  2，System/fsnotify&nbsp; 这是一个监控文件修改的第三方包，已下载到System目录，用户不需再自己安装。该包位于
  <a href="https://github.com/howeyc/fsnotify" target="_blank"> https://github.com/howeyc/fsnotify</a>
</p>
<a href="https://github.com/male110/GoMvc/archive/master.zip">下载GoMvc</a><br/>

<p>有任何问题，可加QQ群：184572648，我基本上每天都在线的</p>
