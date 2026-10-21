<?php
parse_str('a=one+two&list[]=x&list[]=y&user[name]=A%20B', $parsed);
echo $parsed['a'], '|', implode(',', $parsed['list']), '|', $parsed['user']['name'], '|';
echo urldecode('x+y%2Bz'), '|', rawurldecode('x+y%2Bz'), '|';
echo http_build_query(['a b' => 1, 'x' => ['y' => 'z z']], '', '&', PHP_QUERY_RFC3986), '|';
echo hex2bin('4142');
