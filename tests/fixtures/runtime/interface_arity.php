<?php
interface TwoArgs { public function value($first, $second); }
echo "before\n";
class OneArg implements TwoArgs { public function value($first) { return $first; } }
