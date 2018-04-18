# DHT-client-Golang
```  
SQL:
CREATE TABLE `torrent` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `torrent` varchar(255) CHARACTER SET utf8 NOT NULL,
  `filename` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `portnumber` int(11) DEFAULT NULL,
  `init_time` datetime DEFAULT NULL,
  `filesize` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  PRIMARY KEY (`id`,`torrent`),
  KEY `torrent` (`torrent`) USING BTREE,
  KEY `filename` (`filename`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
```  


### todo list ### 
* 使用channel来延时并且批量写入 *  
* 使用command line来传入参数 *  
* 使用配置文件来控制参数 *  