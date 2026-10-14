<?php
interface IncludedContract { public function value(); }
trait IncludedImplementation { public function value() { return 17; } }
