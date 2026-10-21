<?php
function &counter() {
    static $number = 1;
    return $number;
}
$alias =& counter();
$alias = 5;
echo counter(), '|';
$copy = counter();
$copy = 8;
echo counter();
