<?php
$x = 10;
function includeBinding() {
    static $x = 1;
    include __DIR__ . '/../support/review_bind.inc';
    echo $x, "|";
}
includeBinding();
echo $x, "\n";
