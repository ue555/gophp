<?php
interface DnfA {}
interface DnfB {}
class DnfBoth implements DnfA, DnfB {}
function dnfSelect((DnfA&DnfB)|string|null $value): (DnfB&DnfA)|string|null { return $value; }
echo gettype(dnfSelect(new DnfBoth)), '|', gettype(dnfSelect('ok')), '|', gettype(dnfSelect(null)), "\n";
