package chipmunk

/*
Copyright © 2012 Serge Zirukin

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

/*
#include <chipmunk.h>

extern void bbQuery(cpShape *s, void *p);
extern void eachBody_space(cpBody *b, void *p);
extern void eachConstraint_space(cpConstraint *c, void *p);
extern void eachShape_space(cpShape *s, void *p);
extern void nearestPointQuery(cpShape *s, cpFloat distance, cpVect point, void *p);
extern void pointQuery(cpShape *s, void *p);
extern void postStep(cpSpace *space, cpDataPointer key, cpDataPointer data);
extern void segmentQuery(cpShape *s, cpFloat t, cpVect n, void *p);

static inline cpBool space_add_poststep(cpSpace *space, cpDataPointer key, cpDataPointer data) {
  return cpSpaceAddPostStepCallback(space, (void *)postStep, key, data);
}

static inline void space_bb_query(cpSpace *space, cpBB bb, cpLayers layers, cpGroup group, void *f) {
  cpSpaceBBQuery(space, bb, layers, group, bbQuery, f);
}

static inline void space_each_body(cpSpace *space, void *f) {
  cpSpaceEachBody(space, eachBody_space, f);
}

static inline void space_each_constraint(cpSpace *space, void *f) {
  cpSpaceEachConstraint(space, eachConstraint_space, f);
}

static inline void space_each_shape(cpSpace *space, void *f) {
  cpSpaceEachShape(space, eachShape_space, f);
}

static inline void space_nearest_point_query(cpSpace *space, cpVect point, cpFloat maxDistance,
                                             cpLayers layers, cpGroup group, void *f) {
  cpSpaceNearestPointQuery(space, point, maxDistance, layers, group, nearestPointQuery, f);
}

static inline void space_point_query(cpSpace *s, cpVect point, cpLayers layers, cpGroup group, void *p) {
  cpSpacePointQuery(s, point, layers, group, pointQuery, p);
}

static inline void space_segment_query(cpSpace *space,
                                       cpVect   start,
                                       cpVect   end,
                                       cpLayers layers,
                                       cpGroup  group,
                                       void    *f) {
  cpSpaceSegmentQuery(space, start, end, layers, group, segmentQuery, f);
}
*/
import "C"

import (
  "fmt"
  . "unsafe"
  "unsafe"
)

////////////////////////////////////////////////////////////////////////////////

// BBQuery is a rectangle query callback function type.
type BBQuery func(s Shape)

// NearestPointQuery is a callback function type for NearestPointQuery function.
type NearestPointQuery func(s Shape, distance float64, point Vect)

// PointQuery is a callback function type for PointQuery function.
type PointQuery func(s Shape)

// SegmentQuery is a query callback function type.
type SegmentQuery func(s Shape, t float64, n Vect)

type Space interface {
  ActivateShapesTouchingShape(Shape)
  Add(SpaceObject) SpaceObject
  AddBody(Body) Body
  AddConstraint(Constraint) Constraint
  AddPostStepCallback(func(Space, interface{}), interface{}) bool
  AddShape(Shape) Shape
  AddStaticShape(Shape) Shape
  BBQuery(BB, Layers, Group, BBQuery)
  CollisionBias() float64
  CollisionPersistence() Timestamp
  CollisionSlop() float64
  Contains(SpaceObject) bool
  CurrentTimeStep() float64
  Damping() float64
  Each(interface{})
  EachBody(func(Body))
  EachConstraint(func(Constraint))
  EachShape(func(Shape))
  EnableContactGraph() bool
  Free()
  FreeChildren()
  Gravity() Vect
  IdleSpeedThreshold() float64
  IsLocked() bool
  Iterations() int
  NearestPointQuery(Vect, float64, Layers, Group, NearestPointQuery)
  PointQuery(Vect, Layers, Group, PointQuery)
  PointQueryFirst(Vect, Layers, Group) Shape
  ReindexShape(Shape)
  ReindexShapesForBody(Body)
  ReindexStatic()
  Remove(SpaceObject)
  RemoveBody(Body)
  RemoveCollisionHandler(CollisionType, CollisionType)
  RemoveConstraint(Constraint)
  RemoveShape(Shape)
  RemoveStaticShape(Shape)
  SegmentQuery(Vect, Vect, Layers, Group, SegmentQuery)
  SetCollisionBias(float64)
  SetCollisionPersistence(Timestamp)
  SetCollisionSlop(float64)
  SetDamping(float64)
  SetEnableContactGraph(bool)
  SetGravity(Vect)
  SetIdleSpeedThreshold(float64)
  SetIterations(i int)
  SetSleepTimeThreshold(float64)
  SetUserData(interface{})
  SleepTimeThreshold() float64
  StaticBody() Body
  Step(float64)
  UseSpatialHash(float64, int)
  UserData() interface{}
  c() *C.cpSpace
  freeObject(SpaceObject)
}

// SpaceObject is an interface every space object must implement.
type SpaceObject interface {
  Free()
  addToSpace(Space)
  containedInSpace(Space) bool
  removeFromSpace(Space)
}

// Spаce is a basic unit of simulation in Chipmunk.
type Spаce uintptr

////////////////////////////////////////////////////////////////////////////////

var (
  postStepCallbackMap = make(map[Spаce]*map[interface{}]func(Space, interface{}))
)

////////////////////////////////////////////////////////////////////////////////

// ActivateShapesTouchingShape activates body (calls Activate()) of any shape
// that overlaps the given shape.
func (s Spаce) ActivateShapesTouchingShape(sh Shape) {
  C.cpSpaceActivateShapesTouchingShape(s.c(), sh.c())
}

// Add adds an object from space.
func (s Spаce) Add(obj SpaceObject) SpaceObject {
  obj.addToSpace(s)
  return obj
}

// AddBody adds a rigid body to the simulation.
func (s Spаce) AddBody(b Body) Body {
  return cpBody(C.cpSpaceAddBody(s.c(), b.c()))
}

// AddConstraint adds a constraint to the simulation.
func (s Spаce) AddConstraint(c Constraint) Constraint {
  return cpConstraint(C.cpSpaceAddConstraint(s.c(), c.c()))
}

// AddPostStepCallback schedules a post-step callback to be called when Space.Step() finishes.
// You can only register one callback per unique value for key.
// Returns true only if the key has never been scheduled before.
func (s Spаce) AddPostStepCallback(f func(Space, interface{}), key interface{}) bool {
  (*postStepCallbackMap[s])[key] = f
  return cpBool(C.space_add_poststep(s.c(), dataToC(key), dataToC(f)))
}

// AddShape adds a collision shape to the simulation.
// If the shape is attached to a static body, it will be added as a static shape.
func (s Spаce) AddShape(sh Shape) Shape {
  return cpShape(C.cpSpaceAddShape(s.c(), sh.c()))
}

// AddStaticShape explicity adds a shape as a static shape to the simulation.
func (s Spаce) AddStaticShape(sh Shape) Shape {
  return cpShape(C.cpSpaceAddStaticShape(s.c(), sh.c()))
}

// BBQuery performs a fast rectangle query on the space calling a callback
// function for each shape found.
// Only the shape's bounding boxes are checked for overlap, not their full shape.
func (s Spаce) BBQuery(bb BB, layers Layers, group Group, f BBQuery) {
  C.space_bb_query(s.c(), bb.c(), layers.c(), group.c(), Pointer(&f))
}

// CollisionBias returns the speed of how fast overlapping shapes are pushed apart.
func (s Spаce) CollisionBias() float64 {
  return float64(C.cpSpaceGetCollisionBias(s.c()))
}

// CollisionSlop returns amount of encouraged penetration between colliding shapes.
func (s Spаce) CollisionSlop() float64 {
  return float64(C.cpSpaceGetCollisionSlop(s.c()))
}

// CollisionPersistence returns the number of frames that contact information should persist.
func (s Spаce) CollisionPersistence() Timestamp {
  return Timestamp(C.cpSpaceGetCollisionPersistence(s.c()))
}

// Contains tests if a collision shape, rigid body or a constraint has been added to the space.
func (s Spаce) Contains(obj SpaceObject) bool {
  return obj.containedInSpace(s)
}

// CurrentTimeStep returns the current (if you are in a callback from SpaceStep())
// or most recent (outside of a SpaceStep() call) timestep.
func (s Spаce) CurrentTimeStep() float64 {
  return float64(C.cpSpaceGetCurrentTimeStep(s.c()))
}

// Damping returns the damping rate expressed as the fraction of velocity bodies retain each second.
func (s Spаce) Damping() float64 {
  return float64(C.cpSpaceGetDamping(s.c()))
}

// Each calls a callback function on each object of specific type (according to iterator) in the space.
func (s Spаce) Each(iter interface{}) {
  switch iter.(type) {
  case func(Body):
    s.EachBody(iter.(func(Body)))
  case func(Constraint):
    s.EachConstraint(iter.(func(Constraint)))
  case func(Shape):
    s.EachShape(iter.(func(Shape)))
  default:
    panic("invalid iterator in Each()")
  }
}

// EachBody calls a callback function on each body in the space.
func (s Spаce) EachBody(iter func(Body)) {
  p := Pointer(&iter)
  C.space_each_body(s.c(), p)
}

// EachConstraint calls a callback function on each constraint in the space.
func (s Spаce) EachConstraint(iter func(Constraint)) {
  p := Pointer(&iter)
  C.space_each_constraint(s.c(), p)
}

// EachShape calls a callback function on each shape in the space.
func (s Spаce) EachShape(iter func(Shape)) {
  p := Pointer(&iter)
  C.space_each_shape(s.c(), p)
}

// EnableContactGraph returns true if rebuild of the contact graph during each step is enabled.
func (s Spаce) EnableContactGraph() bool {
  return 0 != int(C.cpSpaceGetEnableContactGraph(s.c()))
}

// Free removes a space.
func (s Spаce) Free() {
  delete(postStepCallbackMap, s)
  C.cpSpaceFree(s.c())
}

// FreeChildren frees all bodies, constraints and shapes in the space.
func (s Spаce) FreeChildren() {
  remove := func(s Space, obj interface{}) {
    s.Remove(obj.(SpaceObject))
    s.freeObject(obj.(SpaceObject))
  }

  s.EachShape(func(shape Shape) {
    s.AddPostStepCallback(remove, shape)
  })

  s.EachConstraint(func(constraint Constraint) {
    s.AddPostStepCallback(remove, constraint)
  })

  s.EachBody(func(body Body) {
    s.AddPostStepCallback(remove, body)
  })
}

// Gravity returns current gravity used when integrating velocity for rigid bodies.
func (s Spаce) Gravity() Vect {
  return cpVect(C.cpSpaceGetGravity(s.c()))
}

// IdleSpeedThreshold returns speed threshold for a body to be considered idle.
func (s Spаce) IdleSpeedThreshold() float64 {
  return float64(C.cpSpaceGetIdleSpeedThreshold(s.c()))
}

// IsLocked returns true if objects cannot be added/removed inside a callback.
func (s Spаce) IsLocked() bool {
  return cpBool(C.cpSpaceIsLocked(s.c()))
}

// Iterations returns the number of iterations to use in the impulse solver (to solve contacts).
func (s Spаce) Iterations() int {
  return int(C.cpSpaceGetIterations(s.c()))
}

// NearestPointQuery queries the space at a point and calls a callback function for each shape found.
func (s Spаce) NearestPointQuery(
  point Vect,
  maxDistance float64,
  layers Layers,
  group Group,
  f NearestPointQuery) {
  C.space_nearest_point_query(
    s.c(),
    point.c(),
    C.cpFloat(maxDistance),
    layers.c(),
    group.c(),
    Pointer(&f))
}

// PointQuery queries the space at a point and calls a callback function for each shape found.
func (s Spаce) PointQuery(point Vect, layers Layers, group Group, f PointQuery) {
  C.space_point_query(s.c(), point.c(), layers.c(), group.c(), Pointer(&f))
}

// PointQueryFirst queries the space at a point and returns
// the first shape found. Returns nil if no shapes were found.
func (s Spаce) PointQueryFirst(point Vect, layers Layers, group Group) Shape {
  return cpShape(C.cpSpacePointQueryFirst(s.c(), point.c(), layers.c(), group.c()))
}

// Remove removes an object from space.
func (s Spаce) Remove(obj SpaceObject) {
  obj.removeFromSpace(s)
}

// SegmentQuery performs a directed line segment query (like a raycast)
// against the space calling a callback function for each shape intersected.
func (s Spаce) SegmentQuery(start, end Vect, layers Layers, group Group, f SegmentQuery) {
  C.space_segment_query(s.c(), start.c(), end.c(), layers.c(), group.c(), Pointer(&f))
}

// SetGravity sets the gravity to pass to rigid bodies when integrating velocity.
func (s Spаce) SetGravity(g Vect) {
  C.cpSpaceSetGravity(s.c(), g.c())
}

// SetCollisionBias sets the speed of how fast overlapping shapes are pushed apart.
// Expressed as a fraction of the error remaining after each second.
// Defaults to pow(1.0 - 0.1, 60.0) meaning that Chipmunk fixes 10% of overlap each frame at 60Hz.
func (s Spаce) SetCollisionBias(b float64) {
  C.cpSpaceSetCollisionBias(s.c(), C.cpFloat(b))
}

// SetCollisionPersistence sets the number of frames that contact information should persist.
// Defaults to 3. There is probably never a reason to change this value.
func (s Spаce) SetCollisionPersistence(p Timestamp) {
  C.cpSpaceSetCollisionPersistence(s.c(), C.cpTimestamp(p))
}

// SetCollisionSlop sets amount of encouraged penetration between colliding shapes.
// Used to reduce oscillating contacts and keep the collision cache warm.
// Defaults to 0.1. If you have poor simulation quality,
// increase this number as much as possible without allowing visible amounts of overlap.
func (s Spаce) SetCollisionSlop(sl float64) {
  C.cpSpaceSetCollisionSlop(s.c(), C.cpFloat(sl))
}

// SetDamping sets the damping rate expressed as the fraction of velocity bodies retain each second.
// A value of 0.9 would mean that each body's velocity will drop 10% per second.
// The default value is 1.0, meaning no damping is applied.
// Note this damping value is different than those of DampedSpring and DampedRotarySpring.
func (s Spаce) SetDamping(d float64) {
  C.cpSpaceSetDamping(s.c(), C.cpFloat(d))
}

// SetEnableContactGraph enables a rebuild of the contact graph during each step.
// Must be enabled to use the EachArbiter() method of Body.
// Disabled by default for a small performance boost.
// Enabled implicitly when the sleeping feature is enabled.
func (s Spаce) SetEnableContactGraph(cg bool) {
  C.cpSpaceSetEnableContactGraph(s.c(), boolToC(cg))
}

// SetIdleSpeedThreshold sets the speed threshold for a body to be considered idle.
// The default value of 0.0 means to let the space guess a good threshold based on gravity.
func (s Spаce) SetIdleSpeedThreshold(t float64) {
  C.cpSpaceSetIdleSpeedThreshold(s.c(), C.cpFloat(t))
}

// SetIterations sets the number of iterations to use in the impulse solver to solve contacts.
func (s Spаce) SetIterations(i int) {
  C.cpSpaceSetIterations(s.c(), C.int(i))
}

// SetSleepTimeThreshold sets the time a group of bodies must remain idle in order to fall asleep.
// Enabling sleeping also implicitly enables the the contact graph.
// The default value of math.Inf(1) disables the sleeping algorithm.
func (s Spаce) SetSleepTimeThreshold(t float64) {
  C.cpSpaceSetSleepTimeThreshold(s.c(), C.cpFloat(t))
}

// SetUserData sets user definable data pointer.
// Generally this points to your game's controller or game state
// so you can access it when given a Space reference in a callback.
func (s Spаce) SetUserData(data interface{}) {
  C.cpSpaceSetUserData(s.c(), dataToC(data))
}

// SleepTimeThreshold returns the time a groups of bodies must remain idle in order to "fall asleep".
func (s Spаce) SleepTimeThreshold() float64 {
  return float64(C.cpSpaceGetSleepTimeThreshold(s.c()))
}

// SpaceNew creates a new space.
func SpaceNew() Space {
  s := Spаce(Pointer(C.cpSpaceNew()))
  m := make(map[interface{}]func(Space, interface{}))
  postStepCallbackMap[s] = &m
  return s
}

// String converts a space to a human-readable string.
func (s Spаce) String() string {
  return fmt.Sprintf("(Space)%+v", s)
}

// StaticBody returns a dedicated static body for the space.
// You don't have to use it, but because it's memory is managed automatically with the space
// it's very convenient.
// You can set its user data pointer to something helpful if you want for callbacks.
func (s Spаce) StaticBody() Body {
  return cpBody(C.cpSpaceGetStaticBody(s.c()))
}

// Step makes the space step forward in time by dt seconds.
func (s Spаce) Step(dt float64) {
  C.cpSpaceStep(s.c(), C.cpFloat(dt))
}

// ReindexShape updates the collision detection data for a specific shape in the space.
func (s Spаce) ReindexShape(sh Shape) {
  C.cpSpaceReindexShape(s.c(), sh.c())
}

// ReindexShapesForBody updates the collision detection data for all shapes attached to a body.
func (s Spаce) ReindexShapesForBody(b Body) {
  C.cpSpaceReindexShapesForBody(s.c(), b.c())
}

// ReindexStatic updates the collision detection info for the static shape in the space.
func (s Spаce) ReindexStatic() {
  C.cpSpaceReindexStatic(s.c())
}

// RemoveBody removes a rigid body from the simulation.
func (s Spаce) RemoveBody(b Body) {
  C.cpSpaceRemoveBody(s.c(), b.c())
}

// RemoveCollisionHandler unsets a collision handler.
func (s Spаce) RemoveCollisionHandler(a CollisionType, b CollisionType) {
  C.cpSpaceRemoveCollisionHandler(s.c(), C.cpCollisionType(a), C.cpCollisionType(b))
}

// RemoveConstraint removes a constraint from the simulation.
func (s Spаce) RemoveConstraint(c Constraint) {
  C.cpSpaceRemoveConstraint(s.c(), c.c())
}

// RemoveShape removes a collision shape from the simulation.
func (s Spаce) RemoveShape(sh Shape) {
  C.cpSpaceRemoveShape(s.c(), sh.c())
}

// RemoveStaticShape removes a collision shape added using AddStaticShape() from the simulation.
func (s Spаce) RemoveStaticShape(sh Shape) {
  C.cpSpaceRemoveStaticShape(s.c(), sh.c())
}

// UserData returns user defined data.
func (s Spаce) UserData() interface{} {
  return cpData(C.cpSpaceGetUserData(s.c()))
}

// UseSpatialHash switches the space to use a spatial has as it's spatial index.
func (s Spаce) UseSpatialHash(dim float64, count int) {
  C.cpSpaceUseSpatialHash(s.c(), C.cpFloat(dim), C.int(count))
}

//export bbQuery
func bbQuery(s *C.cpShape, p unsafe.Pointer) {
  f := *(*BBQuery)(p)
  f(cpShape(s))
}

// c converts Space to c.cpSpace pointer.
func (s Spаce) c() *C.cpSpace {
  return (*C.cpSpace)(Pointer(s))
}

// cpSpace converts C.cpSpace pointer to Space.
func cpSpace(s *C.cpSpace) Space {
  if s != nil {
    return Spаce(Pointer(s))
  }

  return nil
}

//export eachBody_space
func eachBody_space(b *C.cpBody, p unsafe.Pointer) {
  f := *(*func(Body))(p)
  f(cpBody(b))
}

//export eachConstraint_space
func eachConstraint_space(c *C.cpConstraint, p unsafe.Pointer) {
  f := *(*func(Constraint))(p)
  f(cpConstraint(c))
}

//export eachShape_space
func eachShape_space(sh *C.cpShape, p unsafe.Pointer) {
  f := *(*func(Shape))(p)
  f(cpShape(sh))
}

// freeObject frees an object.
func (s Spаce) freeObject(obj SpaceObject) {
  obj.Free()
}

//export nearestPointQuery
func nearestPointQuery(s *C.cpShape, distance C.cpFloat, point C.cpVect, p unsafe.Pointer) {
  f := *(*NearestPointQuery)(p)
  f(cpShape(s), float64(distance), cpVect(point))
}

//export pointQuery
func pointQuery(s *C.cpShape, p unsafe.Pointer) {
  f := *(*PointQuery)(p)
  f(cpShape(s))
}

//export postStep
func postStep(s *C.cpSpace, p, data C.cpDataPointer) {
  key := cpData(p)
  f := cpData(data).(func(Space, interface{}))

  // execute callback
  f(cpSpace(s), key)
  // remove from map
  delete(*postStepCallbackMap[Spаce(Pointer(s))], key)
}

//export segmentQuery
func segmentQuery(s *C.cpShape, t C.cpFloat, n C.cpVect, p unsafe.Pointer) {
  f := *(*SegmentQuery)(p)
  f(cpShape(s), float64(t), cpVect(n))
}

// Local Variables:
// indent-tabs-mode: nil
// tab-width: 2
// End:
// ex: set tabstop=2 shiftwidth=2 expandtab:
