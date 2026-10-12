<?php
class Dynamic {}
$o = new Dynamic;
$o->x = 1;
$o->x = 2;
echo $o->x, "\n";
