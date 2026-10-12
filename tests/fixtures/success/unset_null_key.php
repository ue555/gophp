<?php
$a = ["" => 1]; unset($a[null]); echo count($a), "\n";
