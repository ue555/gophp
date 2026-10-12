<?php
function gophp_isset_order_key($n) { echo $n; return 0; }
$a = [7]; $n = null;
echo ':', isset($a[gophp_isset_order_key(1)], $n, $a[gophp_isset_order_key(2)]) + 0, "\n";
echo ':', isset($a[gophp_isset_order_key(3)], $a[gophp_isset_order_key(4)]) + 0, "\n";
echo ':', isset($missing[gophp_isset_order_key(5)][gophp_isset_order_key(6)]) + 0, "\n";
