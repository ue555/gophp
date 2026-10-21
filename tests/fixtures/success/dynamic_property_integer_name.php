<?php
$box = new stdClass;
$name = 7;
$box->{$name} = 'seven';
echo $box->$name, '|', $box->{'7'}, "\n";
