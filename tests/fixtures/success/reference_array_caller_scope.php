<?php
function change($copy) {
    $copy[0] = 9;
}

function outer() {
    $a = [1];
    $alias =& $a[0];
    change($a);
    echo $a[0], '|', $alias;
}

outer();
