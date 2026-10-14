<?php
$x = 10;
function reverseBindingOrder() {
    global $x;
    static $x = 1;
    $x++;
    echo $x, "|";
}
reverseBindingOrder();
echo $x, "\n";
