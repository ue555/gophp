<?php
$original = [1];
$copy = $original;
$alias =& $original;
$copy[0] = 2;
$alias[0] = 3;
echo $original[0], '|', $copy[0], '|', $alias[0];
