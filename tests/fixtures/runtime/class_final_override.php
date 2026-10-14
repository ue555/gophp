<?php
echo "before\n";
class FinalMethodParent { final public function value() {} }
class FinalMethodChild extends FinalMethodParent { public function value() {} }
