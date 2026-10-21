<?php
class DynamicMissing {}
$class = 'DynamicMissing'; $method = 'gone';
echo $class::$method();
