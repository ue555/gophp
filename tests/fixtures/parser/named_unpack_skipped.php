<?php
function f($a, $b) {}
echo "before\n";
if (false) {
    f(a: 1, ...['b' => 2]);
}
