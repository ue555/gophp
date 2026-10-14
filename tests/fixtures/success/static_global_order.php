<?php
$x = 10;
function bindingOrder() {
    static $x = 1;
    global $x;
    $x++;
    echo $x, "|";
}
bindingOrder();
echo $x, "\n";
