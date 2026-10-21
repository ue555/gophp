<?php
echo 'a';
ob_start();
echo 'b';
ob_start();
echo 'c';
$inner = ob_get_clean();
echo $inner;
$outer = ob_get_clean();
echo '|', $outer, '|d';
