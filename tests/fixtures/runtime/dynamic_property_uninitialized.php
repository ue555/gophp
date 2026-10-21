<?php
class UninitializedDynamicBox { public int $value; }
$box = new UninitializedDynamicBox;
$name = 'value';
echo $box->{$name};
