<?php
$o = new stdClass;
echo $o?->missing?->take() ?? "fallback";
