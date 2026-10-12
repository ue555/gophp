<?php
$o = new stdClass;
echo ($o == 1) + 0;
echo ($o < 2.0) + 0;
echo 0 <=> $o, "\n";
