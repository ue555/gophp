<?php
function factorial($n) { if ($n <= 1) return 1; return $n * factorial($n - 1); }
function evenNumber($n) { if ($n === 0) return true; return oddNumber($n - 1); }
function oddNumber($n) { if ($n === 0) return false; return evenNumber($n - 1); }
echo factorial(0), ':', factorial(6), ':', (EVENNUMBER(8) === true) + 0, "\n";
