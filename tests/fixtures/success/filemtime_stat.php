<?php
echo is_int(filemtime(__FILE__)) ? "integer\n":"bad\n";
echo filemtime(__FILE__) > 0 ? "positive\n":"bad\n";
set_error_handler(function($n,$s){echo $s,"\n";return true;});
echo json_encode(filemtime('/gophp-nonexistent-filemtime')), "\n";
try {filemtime([]);} catch(Throwable $e){echo $e->getMessage(),"\n";}
try {filemtime("a\0b");} catch(Throwable $e){echo $e->getMessage(),"\n";}
