<?php
class DynamicPrivate { private static $secret = 7; }
$class = 'DynamicPrivate'; $name = 'secret';
echo $class::${$name};
