<?php
echo "before\n";
class VisibleMethodParent { public function value() {} }
class VisibleMethodChild extends VisibleMethodParent { protected function value() {} }
