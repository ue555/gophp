<?php
class Broken {
    function __construct() { echo 'construct'; $n = 1 / 0; }
}
$o = new Broken;
