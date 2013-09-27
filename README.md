
<a href="https://github.com/male110/GoMvc/archive/master.zip">下载GoMvc</a><br/>
<a href="src/docs/index.html">文档</a> 位于src/docs目录下。

<p>有任何问题，可加群：184572648，我基本上每天都在线的</p>
 
 <a href="#build"> 编译</a><br />
<a href="#config">  配置文件</a><br />
<a href="#route">  路由注册</a>
<p>GoMVC是一个简单，便捷的MVC框架。程序注释全部使用中文，很适合国人使用。文档也很详细。
<a name="build">编译</a>时，需要把GoMvc目录设置为GOPATH.
</p>
<p>
<b><a name="config">配置文件</a></b>
</p>
<div>  
    <p>
        网站的配置文件为web.config，格式为XML，配置项的内容如下：</p>
    <p>
        <b>ShowErrors：</b>是否显示错误信息。true,显示；false,不显示。建义在测试时可以设置为true,发布到正式环境后设置为false。</p>
    <p>
        <b>CookieDomain：</b>Cookies的Domain信息，可用来共享cookie。如domain.com，和sub.domain.com，可以通过把CookieDomain统一设置为domain.com来共享cookies信息</p>
    <p>
        <b>Theme：</b>网站当前使用的主题，在Views目录下，可以有多套网站模板。</p>
    <p>
        <b>LogPath：</b>日志文件的存放位置</p>
    <p>
        <b>LogFileMaxSize：</b>单个日志文件的大小，超过指定大小后将创建一个新的日志文件。</p>
    <p>
        <b>DriverName：</b>数据库的驱动名称。</p>
    <p>
        <b>DataSourceName：</b>数据库的连接字符串。</p>
    <p>
        <b>StaticDir：</b>静态目录,该目录下通常是CSS,JS,图片等静态资源。</p>
    <p>
        <b>StaticFile：</b>静态文件，用来设置单个的静态文件，主要是为了提高灵活性，满足特殊的需求.</p>
    <p>
        <b>SessionType：</b>Session的存放类型,1,文件,2内存,3Mysql数据库,修改需重启才能生效。当配置为3时，需要在数据库中创建一个表，来存放session,创建表的SQL如下：<br />
    </p>
    <pre>CREATE TABLE `session` (
	`session_id` CHAR(32) NULL,
	`session_data` BLOB NULL,
	`lastupdatetime` DATETIME NULL,
	PRIMARY KEY (`session_id`)
)
COLLATE=&#39;utf8_general_ci&#39;;
</pre>
    <p>
        <b>SessionLocation：</b>当SessionType为1时，该项为Session文件的存放路径；SessionType为3时,该项为数据库连接字符串。</p>
    <p>
        <b>SessionTimeOut：</b>Session超时时间，单位分钟</p>
    <p>
        <b>MemFreeInterval：</b>程序中有定时器，定时对Session进行检查，删除超时的Session，该配置项用来设置多久进行一次检查，单位秒，默认值60。</p>
    <p>
        <b>ListenPort：</b>网站的端口号,该配置改后必须重启程序才能生效。</p>
    <p>
        &nbsp;</p>
</div>
<p>
  <b><a name="route">  路由注册</a></b></p>
<p>
    用RouteTable.AddRote来注册路由。其格式如下： 
</p>
<pre>//注册标准路由
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
<pre>RouteTable.AddRote(&amp;RouteItem{
		Name:        &quot;default&quot;,
		Url:         &quot;{controller}/{action}/{id}&quot;,
		Defaults:    map[string]interface{}{&quot;controller&quot;: &quot;home&quot;, &quot;action&quot;: &quot;index&quot;, &quot;id&quot;: 123},
		Constraints: map[string]string{&quot;id&quot;: `^(\d+)$`}})
</pre>
在上面的例子中我们指定了id参数只能是数字，并设置了默认值123。要在Controller中获取该参数值，可以用this.RouteData[&quot;id&quot;]。 
<p>
    因为在Go没有办法反射出包中的所有struct，所以需要手动来注册Controller,格式如下： 
</p>
<pre>import (
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
    &nbsp;</p>
 <a href="https://github.com/male110/GoMvc/archive/master.zip">下载GoMvc</a><br/>
<a href="src/docs/index.html">文档</a> 位于src/docs目录下。
<p>有任何问题，可加群：184572648，我基本上每天都在线的</p>