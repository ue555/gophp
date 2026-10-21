<?php
class PrivateDynamicBox { private int $value = 1; }
$box = new PrivateDynamicBox;
$name = 'value';
echo $box->{$name};
