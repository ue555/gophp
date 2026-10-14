<?php
namespace DynamicNames;
class Box {}
$name = __NAMESPACE__ . '\\Box';
$a = new $name;
$b = new $name();
echo $a instanceof $name, '|', $b instanceof $a, "\n";
