<?php
interface PartialDnfA {}
interface PartialDnfB {}
class OnlyA implements PartialDnfA {}
function dnfNeeds((PartialDnfA&PartialDnfB)|int $x) {}
dnfNeeds(new OnlyA);
