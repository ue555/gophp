<?php
class CreatedDynamicBox {}
$box = new CreatedDynamicBox;
$name = 'extra';
$box->{$name} = 12;
echo $box->$name, "\n";
