<?php
$a = new stdClass; $a->self = $a;
$b = new stdClass; $b->self = $b;
echo 'before';
echo $a == $b;
