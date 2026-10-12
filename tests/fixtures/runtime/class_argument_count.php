<?php
class Required {
    function run($x) { echo 'bad'; }
}
$o = new Required;
$o->run();
