<?php
class Divider {
    function divide($n) { return 1 / $n; }
}
$o = new Divider;
echo 'before';
$o->divide(0);
