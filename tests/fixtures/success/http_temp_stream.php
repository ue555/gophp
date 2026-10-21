<?php
$stream = fopen('php://temp/maxmemory:1048576', 'w+');
fwrite($stream, 'abc');
fseek($stream, 0);
echo fread($stream, 2), '|', stream_get_contents($stream);
fclose($stream);
