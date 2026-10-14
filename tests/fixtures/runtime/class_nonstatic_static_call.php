<?php
class NonStaticCall { public function value() { return 1; } }
echo "before\n";
echo NonStaticCall::value();
