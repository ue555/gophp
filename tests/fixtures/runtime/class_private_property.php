<?php
class PrivateProperty { private $x = 1; }
$o = new PrivateProperty;
echo $o->x;
