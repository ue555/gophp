<?php
class RefDynamic { public $value = 1; }
$box = new RefDynamic;
$events = '';
function refReceiver() { global $box, $events; $events .= 'R'; return $box; }
function refName() { global $events; $events .= 'N'; return 'value'; }
$alias =& refReceiver()->{refName()};
$alias = 6;
echo $events, '|', $box->value;
