// noinspection DuplicatedCode

import * as jo from "jo";
import {int, float, boolean} from "jo";
import {assert} from "jo";
import {ContactFilter} from "org/jbox2d/callbacks/ContactFilter.js";
import {ContactListener} from "org/jbox2d/callbacks/ContactListener.js";
import {DebugDraw} from "org/jbox2d/callbacks/DebugDraw.js";
import {DestructionListener} from "org/jbox2d/callbacks/DestructionListener.js";
import {ParticleDestructionListener} from "org/jbox2d/callbacks/ParticleDestructionListener.js";
import {ParticleQueryCallback} from "org/jbox2d/callbacks/ParticleQueryCallback.js";
import {ParticleRaycastCallback} from "org/jbox2d/callbacks/ParticleRaycastCallback.js";
import {QueryCallback} from "org/jbox2d/callbacks/QueryCallback.js";
import {RayCastCallback} from "org/jbox2d/callbacks/RayCastCallback.js";
import {TreeCallback} from "org/jbox2d/callbacks/TreeCallback.js";
import {TreeRayCastCallback} from "org/jbox2d/callbacks/TreeRayCastCallback.js";
import {AABB} from "org/jbox2d/collision/AABB.js";
import {RayCastInput} from "org/jbox2d/collision/RayCastInput.js";
import {RayCastOutput} from "org/jbox2d/collision/RayCastOutput.js";
import {TOIInput} from "org/jbox2d/collision/TimeOfImpact/TOIInput.js";
import {TOIOutput} from "org/jbox2d/collision/TimeOfImpact/TOIOutput.js";
import {TOIOutputState} from "org/jbox2d/collision/TimeOfImpact/TOIOutputState.js";
import {BroadPhase} from "org/jbox2d/collision/broadphase/BroadPhase.js";
import {BroadPhaseStrategy} from "org/jbox2d/collision/broadphase/BroadPhaseStrategy.js";
import {DefaultBroadPhaseBuffer} from "org/jbox2d/collision/broadphase/DefaultBroadPhaseBuffer.js";
import {DynamicTree} from "org/jbox2d/collision/broadphase/DynamicTree.js";
import {ChainShape} from "org/jbox2d/collision/shapes/ChainShape.js";
import {CircleShape} from "org/jbox2d/collision/shapes/CircleShape.js";
import {EdgeShape} from "org/jbox2d/collision/shapes/EdgeShape.js";
import {PolygonShape} from "org/jbox2d/collision/shapes/PolygonShape.js";
import {Shape} from "org/jbox2d/collision/shapes/Shape.js";
import {ShapeType} from "org/jbox2d/collision/shapes/ShapeType.js";
import {Color3f} from "org/jbox2d/common/Color3f.js";
import {MathUtils} from "org/jbox2d/common/MathUtils.js";
import {Settings} from "org/jbox2d/common/Settings.js";
import {Sweep} from "org/jbox2d/common/Sweep.js";
import {Timer} from "org/jbox2d/common/Timer.js";
import {Transform} from "org/jbox2d/common/Transform.js";
import {Vec2} from "org/jbox2d/common/Vec2.js";
import {Contact} from "org/jbox2d/dynamics/contacts/Contact.js";
import {ContactEdge} from "org/jbox2d/dynamics/contacts/ContactEdge.js";
import {ContactRegister} from "org/jbox2d/dynamics/contacts/ContactRegister.js";
import {joints} from "org/jbox2d/dynamics/joints.js";
import {ParticleBodyContact} from "org/jbox2d/particle/ParticleBodyContact.js";
import {ParticleColor} from "org/jbox2d/particle/ParticleColor.js";
import {ParticleContact} from "org/jbox2d/particle/ParticleContact.js";
import {ParticleDef} from "org/jbox2d/particle/ParticleDef.js";
import {ParticleGroup} from "org/jbox2d/particle/ParticleGroup.js";
import {ParticleGroupDef} from "org/jbox2d/particle/ParticleGroupDef.js";
import {ParticleSystem} from "org/jbox2d/particle/ParticleSystem.js";
import {IDynamicStack} from "org/jbox2d/pooling/IDynamicStack.js";
import {IWorldPool} from "org/jbox2d/pooling/IWorldPool.js";
import {Vec2Array} from "org/jbox2d/pooling/arrays/Vec2Array.js";
import {DefaultWorldPool} from "org/jbox2d/pooling/normal/DefaultWorldPool.js";
import {ContactManager} from 'org/jbox2d/dynamics/ContactManager.js';
import {Body} from 'org/jbox2d/dynamics/Body.js';
import {Profile} from 'org/jbox2d/dynamics/Profile.js';
import {Fixture} from 'org/jbox2d/dynamics/Fixture.js';
import {BodyDef} from 'org/jbox2d/dynamics/BodyDef.js';
import {TimeStep} from 'org/jbox2d/dynamics/TimeStep.js';
import {FixtureProxy} from 'org/jbox2d/dynamics/FixtureProxy.js';
import {Island} from 'org/jbox2d/dynamics/Island.js';
import {BodyType} from 'org/jbox2d/dynamics/BodyType.js';

export class World {

    /**
     * @var {number}
     */
    static WORLD_POOL_SIZE = 100;

    /**
     * @var {number}
     */
    static WORLD_POOL_CONTAINER_SIZE = 10;

    /**
     * @var {number}
     */
    static NEW_FIXTURE = 0x0001;

    /**
     * @var {number}
     */
    static LOCKED = 0x0002;

    /**
     * @var {number}
     */
    static CLEAR_FORCES = 0x0004;

    /**
     * @var {number}
     */
    activeContacts = 0;

    /**
     * @var {number}
     */
    contactPoolCount = 0;

    /**
     * @var {number}
     */
    m_flags;

    /**
     * @var {ContactManager}
     */
    m_contactManager;

    /**
     * @var {Body}
     */
    m_bodyList;

    /**
     * @var {Joint}
     */
    m_jointList;

    /**
     * @var {number}
     */
    m_bodyCount;

    /**
     * @var {number}
     */
    m_jointCount;

    /**
     * @var {Vec2}
     */
    m_gravity = new Vec2();

    /**
     * @var {boolean}
     */
    m_allowSleep;

    /**
     * @var {DestructionListener}
     */
    m_destructionListener;

    /**
     * @var {ParticleDestructionListener}
     */
    m_particleDestructionListener;

    /**
     * @var {DebugDraw}
     */
    m_debugDraw;

    /**
     * @var {IWorldPool}
     */
    pool;

    /**
     * @var {number}
     */
    m_inv_dt0;

    /**
     * @var {boolean}
     */
    m_warmStarting;

    /**
     * @var {boolean}
     */
    m_continuousPhysics;

    /**
     * @var {boolean}
     */
    m_subStepping;

    /**
     * @var {boolean}
     */
    m_stepComplete;

    /**
     * @var {Profile}
     */
    m_profile;

    /**
     * @var {ParticleSystem}
     */
    m_particleSystem;

    /**
     * @var {ContactRegister[][]}
     */
    contactStacks = new ContactRegister[ShapeType.values().length][ShapeType.values().length];

    /**
     * @var {TimeStep}
     */
    step = new TimeStep();

    /**
     * @var {Timer}
     */
    stepTimer = new Timer();

    /**
     * @var {Timer}
     */
    tempTimer = new Timer();

    /**
     * @var {Color3f}
     */
    color = new Color3f();

    /**
     * @var {Transform}
     */
    xf = new Transform();

    /**
     * @var {Vec2}
     */
    cA = new Vec2();

    /**
     * @var {Vec2}
     */
    cB = new Vec2();

    /**
     * @var {Vec2Array}
     */
    avs = new Vec2Array();

    /**
     * @var {WorldQueryWrapper}
     */
    wqwrapper = new WorldQueryWrapper();

    /**
     * @var {WorldRayCastWrapper}
     */
    wrcwrapper = new WorldRayCastWrapper();

    /**
     * @var {RayCastInput}
     */
    input = new RayCastInput();

    /**
     * @var {Island}
     */
    island = new Island();

    /**
     * @var {Body[]}
     */
    stack = new Body[10];

    /**
     * @var {Timer}
     */
    broadphaseTimer = new Timer();

    /**
     * @var {Island}
     */
    toiIsland = new Island();

    /**
     * @var {TOIInput}
     */
    toiInput = new TOIInput();

    /**
     * @var {TOIOutput}
     */
    toiOutput = new TOIOutput();

    /**
     * @var {TimeStep}
     */
    subStep = new TimeStep();

    /**
     * @var {Body[]}
     */
    tempBodies = new Body[2];

    /**
     * @var {Sweep}
     */
    backup1 = new Sweep();

    /**
     * @var {Sweep}
     */
    backup2 = new Sweep();

    /**
     * @var {Integer}
     */
    static LIQUID_INT = 1234598372;

    /**
     * @var {number}
     */
    liquidLength = .12;

    /**
     * @var {number}
     */
    averageLinearVel = -1;

    /**
     * @var {Vec2}
     */
    liquidOffset = new Vec2();

    /**
     * @var {Vec2}
     */
    circCenterMoved = new Vec2();

    /**
     * @var {Color3f}
     */
    liquidColor = new Color3f(.4, .4, 1);

    /**
     * @var {Vec2}
     */
    center = new Vec2();

    /**
     * @var {Vec2}
     */
    axis = new Vec2();

    /**
     * @var {Vec2}
     */
    v1 = new Vec2();

    /**
     * @var {Vec2}
     */
    v2 = new Vec2();

    /**
     * @var {Vec2Array}
     */
    tlvertices = new Vec2Array();

    constructor() {
        const $this = () => {
            switch (true){
                case jo.suitable(arguments, Vec2): {
                    let [gravity] = arguments;
                    $this(gravity, new DefaultWorldPool(World.WORLD_POOL_SIZE, World.WORLD_POOL_CONTAINER_SIZE));
                    break;
                }

                case jo.suitable(arguments, Vec2, IWorldPool): {
                    let [gravity, pool] = arguments;
                    $this(gravity, pool, new DynamicTree());
                    break;
                }

                case jo.suitable(arguments, Vec2, IWorldPool, BroadPhaseStrategy): {
                    let [gravity, pool, strategy] = arguments;
                    $this(gravity, pool, new DefaultBroadPhaseBuffer(strategy));
                    break;
                }

                case jo.suitable(arguments, Vec2, IWorldPool, BroadPhase): {
                    let [gravity, pool, broadPhase] = arguments;
                    this.pool = pool;
                    this.m_destructionListener = null;
                    this.m_debugDraw = null;
                    this.m_bodyList = null;
                    this.m_jointList = null;
                    this.m_bodyCount = 0;
                    this.m_jointCount = 0;
                    this.m_warmStarting = true;
                    this.m_continuousPhysics = true;
                    this.m_subStepping = false;
                    this.m_stepComplete = true;
                    this.m_allowSleep = true;
                    m_gravity.set(gravity);
                    this.m_flags = World.CLEAR_FORCES;
                    this.m_inv_dt0 = 0;
                    this.m_contactManager = new ContactManager(this, broadPhase);
                    this.m_profile = new Profile();
                    this.m_particleSystem = new ParticleSystem(this);
                    this.initializeRegisters();
                    break;
                }

                // default:
                //     super(...arguments);
                //     break;
            }
        }

        $this(...arguments)
    }

    destroyParticlesInShape() {
        if (jo.suitable(arguments, Shape, Transform)) return this.destroyParticlesInShape$ShapeTransform(...arguments);
        if (jo.suitable(arguments, Shape, Transform, boolean)) return this.destroyParticlesInShape$ShapeTransformboolean(...arguments);

        return super.destroyParticlesInShape(...arguments);
    }

    destroyParticle() {
        if (jo.suitable(arguments, int)) return this.destroyParticle$int(...arguments);
        if (jo.suitable(arguments, int, boolean)) return this.destroyParticle$intboolean(...arguments);

        return super.destroyParticle(...arguments);
    }

    queryAABB() {
        if (jo.suitable(arguments, QueryCallback, AABB)) return this.queryAABB$QueryCallbackAABB(...arguments);
        if (jo.suitable(arguments, QueryCallback, ParticleQueryCallback, AABB)) return this.queryAABB$QueryCallbackParticleQueryCallbackAABB(...arguments);
        if (jo.suitable(arguments, ParticleQueryCallback, AABB)) return this.queryAABB$ParticleQueryCallbackAABB(...arguments);

        return super.queryAABB(...arguments);
    }

    raycast() {
        if (jo.suitable(arguments, RayCastCallback, Vec2, Vec2)) return this.raycast$RayCastCallbackVec2Vec2(...arguments);
        if (jo.suitable(arguments, RayCastCallback, ParticleRaycastCallback, Vec2, Vec2)) return this.raycast$RayCastCallbackParticleRaycastCallbackVec2Vec2(...arguments);
        if (jo.suitable(arguments, ParticleRaycastCallback, Vec2, Vec2)) return this.raycast$ParticleRaycastCallbackVec2Vec2(...arguments);

        return super.raycast(...arguments);
    }

    destroyParticlesInGroup() {
        if (jo.suitable(arguments, ParticleGroup, boolean)) return this.destroyParticlesInGroup$ParticleGroupboolean(...arguments);
        if (jo.suitable(arguments, ParticleGroup)) return this.destroyParticlesInGroup$ParticleGroup(...arguments);

        return super.destroyParticlesInGroup(...arguments);
    }


    setAllowSleep(flag) {
        if (flag == this.m_allowSleep) {
            return;
        }
        this.m_allowSleep = flag;
        if (this.m_allowSleep == false) {
            for (let b = m_bodyList;
                 b != null; b = b.m_next) {
                b.setAwake(true);
            }
        }
    }

    setSubStepping(subStepping) {
        this.m_subStepping = subStepping;
    }

    isSubStepping() {
        return this.m_subStepping;
    }

    isAllowSleep() {
        return this.m_allowSleep;
    }

    addType(creator, type1, type2) {
        let register = new ContactRegister();
        register.creator = creator;
        register.primary = true;
        this.contactStacks[type1.ordinal()][type2.ordinal()] = register;
        if (type1 != type2) {
            let register2 = new ContactRegister();
            register2.creator = creator;
            register2.primary = false;
            this.contactStacks[type2.ordinal()][type1.ordinal()] = register2;
        }
    }

    initializeRegisters() {
        this.addType(pool.getCircleContactStack(), ShapeType.CIRCLE, ShapeType.CIRCLE);
        this.addType(pool.getPolyCircleContactStack(), ShapeType.POLYGON, ShapeType.CIRCLE);
        this.addType(pool.getPolyContactStack(), ShapeType.POLYGON, ShapeType.POLYGON);
        this.addType(pool.getEdgeCircleContactStack(), ShapeType.EDGE, ShapeType.CIRCLE);
        this.addType(pool.getEdgePolyContactStack(), ShapeType.EDGE, ShapeType.POLYGON);
        this.addType(pool.getChainCircleContactStack(), ShapeType.CHAIN, ShapeType.CIRCLE);
        this.addType(pool.getChainPolyContactStack(), ShapeType.CHAIN, ShapeType.POLYGON);
    }

    getDestructionListener() {
        return this.m_destructionListener;
    }

    getParticleDestructionListener() {
        return this.m_particleDestructionListener;
    }

    setParticleDestructionListener(listener) {
        this.m_particleDestructionListener = listener;
    }

    popContact(fixtureA, indexA, fixtureB, indexB) {
        const type1 = fixtureA.getType();
        const type2 = fixtureB.getType();
        const reg = this.contactStacks[type1.ordinal()][type2.ordinal()];
        if (reg != null) {
            if (reg.primary) {
                let c = reg.creator.pop();
                c.init(fixtureA, indexA, fixtureB, indexB);
                return c;
            } else {
                let c = reg.creator.pop();
                c.init(fixtureB, indexB, fixtureA, indexA);
                return c;
            }
        } else {
            return null;
        }
    }

    pushContact(contact) {
        let fixtureA = contact.getFixtureA();
        let fixtureB = contact.getFixtureB();
        if (contact.m_manifold.pointCount > 0 && !fixtureA.isSensor() && !fixtureB.isSensor()) {
            fixtureA.getBody().setAwake(true);
            fixtureB.getBody().setAwake(true);
        }
        let type1 = fixtureA.getType();
        let type2 = fixtureB.getType();
        let creator = this.contactStacks[type1.ordinal()][type2.ordinal()].creator;
        creator.push(contact);
    }

    getPool() {
        return this.pool;
    }

    setDestructionListener(listener) {
        this.m_destructionListener = listener;
    }

    setContactFilter(filter) {
        this.m_contactManager.m_contactFilter = filter;
    }

    setContactListener(listener) {
        this.m_contactManager.m_contactListener = listener;
    }

    setDebugDraw(debugDraw) {
        this.m_debugDraw = debugDraw;
    }

    createBody(def) {
        assert(this.isLocked() == false);
        if (this.isLocked()) {
            return null;
        }// TODO djm pooling
        let b = new Body(def, this);
        // add to world doubly linked list
        b.m_prev = null;
        b.m_next = this.m_bodyList;
        if (this.m_bodyList != null) {
            this.m_bodyList.m_prev = b;
        }
        this.m_bodyList = b;
        ++this.m_bodyCount;
        return b;
    }

    destroyBody(body) {
        assert(this.m_bodyCount > 0);
        assert(this.isLocked() == false);
        if (this.isLocked()) {
            return;
        }// Delete the attached joints.
        let je = body.m_jointList;
        while (je != null) {
            let je0 = je;
            je = je.next;
            if (this.m_destructionListener != null) {
                m_destructionListener.sayGoodbye(je0.joint);
            }
            this.destroyJoint(je0.joint);
            body.m_jointList = je;
        }
        body.m_jointList = null;
        // Delete the attached contacts.
        let ce = body.m_contactList;
        while (ce != null) {
            let ce0 = ce;
            ce = ce.next;
            m_contactManager.destroy(ce0.contact);
        }
        body.m_contactList = null;
        let f = body.m_fixtureList;
        while (f != null) {
            let f0 = f;
            f = f.m_next;
            if (this.m_destructionListener != null) {
                m_destructionListener.sayGoodbye(f0);
            }
            f0.destroyProxies(this.m_contactManager.m_broadPhase);
            f0.destroy();
            // TODO djm recycle fixtures (here or in that destroy method)
            body.m_fixtureList = f;
            body.m_fixtureCount = 1;
        }
        body.m_fixtureList = null;
        body.m_fixtureCount = 0;
        // Remove world body list.
        if (body.m_prev != null) {
            body.m_prev.m_next = body.m_next;
        }
        if (body.m_next != null) {
            body.m_next.m_prev = body.m_prev;
        }
        if (body == this.m_bodyList) {
            this.m_bodyList = body.m_next;
        }
        --this.m_bodyCount;
        // TODO djm recycle body
    }

    createJoint(def) {
        assert(this.isLocked() == false);
        if (this.isLocked()) {
            return null;
        }
        let j = Joint.create(this, def);
        // Connect to the world list.
        j.m_prev = null;
        j.m_next = this.m_jointList;
        if (this.m_jointList != null) {
            this.m_jointList.m_prev = j;
        }
        this.m_jointList = j;
        ++this.m_jointCount;
        // Connect to the bodies' doubly linked lists.
        j.m_edgeA.joint = j;
        j.m_edgeA.other = j.getBodyB();
        j.m_edgeA.prev = null;
        j.m_edgeA.next = j.getBodyA().m_jointList;
        if (j.getBodyA().m_jointList != null) {
            j.getBodyA().m_jointList.prev = j.m_edgeA;
        }
        j.getBodyA().m_jointList = j.m_edgeA;
        j.m_edgeB.joint = j;
        j.m_edgeB.other = j.getBodyA();
        j.m_edgeB.prev = null;
        j.m_edgeB.next = j.getBodyB().m_jointList;
        if (j.getBodyB().m_jointList != null) {
            j.getBodyB().m_jointList.prev = j.m_edgeB;
        }
        j.getBodyB().m_jointList = j.m_edgeB;
        let bodyA = def.bodyA;
        let bodyB = def.bodyB;
        // If the joint prevents collisions, then flag any contacts for filtering.
        if (def.collideConnected == false) {
            let edge = bodyB.getContactList();
            while (edge != null) {
                if (edge.other == bodyA) {
                    // Flag the contact for filtering at the next time step (where either
                    // body is awake).
                    edge.contact.flagForFiltering();
                }
                edge = edge.next;
            }
        }// Note: creating a joint doesn't wake the bodies.
        return j;
    }

    destroyJoint(j) {
        assert(this.isLocked() == false);
        if (this.isLocked()) {
            return;
        }
        let collideConnected = j.getCollideConnected();
        // Remove from the doubly linked list.
        if (j.m_prev != null) {
            j.m_prev.m_next = j.m_next;
        }
        if (j.m_next != null) {
            j.m_next.m_prev = j.m_prev;
        }
        if (j == this.m_jointList) {
            this.m_jointList = j.m_next;
        }// Disconnect from island graph.
        let bodyA = j.getBodyA();
        let bodyB = j.getBodyB();
        // Wake up connected bodies.
        bodyA.setAwake(true);
        bodyB.setAwake(true);
        // Remove from body 1.
        if (j.m_edgeA.prev != null) {
            j.m_edgeA.prev.next = j.m_edgeA.next;
        }
        if (j.m_edgeA.next != null) {
            j.m_edgeA.next.prev = j.m_edgeA.prev;
        }
        if (j.m_edgeA == bodyA.m_jointList) {
            bodyA.m_jointList = j.m_edgeA.next;
        }
        j.m_edgeA.prev = null;
        j.m_edgeA.next = null;
        // Remove from body 2
        if (j.m_edgeB.prev != null) {
            j.m_edgeB.prev.next = j.m_edgeB.next;
        }
        if (j.m_edgeB.next != null) {
            j.m_edgeB.next.prev = j.m_edgeB.prev;
        }
        if (j.m_edgeB == bodyB.m_jointList) {
            bodyB.m_jointList = j.m_edgeB.next;
        }
        j.m_edgeB.prev = null;
        j.m_edgeB.next = null;
        Joint.destroy(j);
        assert(this.m_jointCount > 0);
        --this.m_jointCount;
        // If the joint prevents collisions, then flag any contacts for filtering.
        if (collideConnected == false) {
            let edge = bodyB.getContactList();
            while (edge != null) {
                if (edge.other == bodyA) {
                    // Flag the contact for filtering at the next time step (where either
                    // body is awake).
                    edge.contact.flagForFiltering();
                }
                edge = edge.next;
            }
        }
    }

    step(dt, velocityIterations, positionIterations) {
        stepTimer.reset();
        tempTimer.reset();
        // log.debug("Starting step");
        // If new fixtures were added, we need to find the new contacts.
        if ((this.m_flags & World.NEW_FIXTURE) == World.NEW_FIXTURE) {
            // log.debug("There's a new fixture, lets look for new contacts");
            m_contactManager.findNewContacts();
            this.m_flags = ~World.NEW_FIXTURE;
        }
        this.m_flags = World.LOCKED;
        this.step.dt = dt;
        this.step.velocityIterations = velocityIterations;
        this.step.positionIterations = positionIterations;
        if (dt > 0.0) {
            this.step.inv_dt = 1.0 / dt;
        } else {
            this.step.inv_dt = 0.0;
        }
        this.step.dtRatio = this.m_inv_dt0 * dt;
        this.step.warmStarting = this.m_warmStarting;
        this.m_profile.stepInit.record(tempTimer.getMilliseconds());
        // Update contacts. This is where some contacts are destroyed.
        tempTimer.reset();
        m_contactManager.collide();
        this.m_profile.collide.record(tempTimer.getMilliseconds());
        // Integrate velocities, solve velocity constraints, and integrate positions.
        if (this.m_stepComplete && this.step.dt > 0.0) {
            tempTimer.reset();
            m_particleSystem.solve(this.step);
            // Particle Simulation
            this.m_profile.solveParticleSystem.record(tempTimer.getMilliseconds());
            tempTimer.reset();
            this.solve(this.step);
            this.m_profile.solve.record(tempTimer.getMilliseconds());
        }// Handle TOI events.
        if (this.m_continuousPhysics && this.step.dt > 0.0) {
            tempTimer.reset();
            this.solveTOI(this.step);
            this.m_profile.solveTOI.record(tempTimer.getMilliseconds());
        }
        if (this.step.dt > 0.0) {
            this.m_inv_dt0 = this.step.inv_dt;
        }
        if ((this.m_flags & World.CLEAR_FORCES) == World.CLEAR_FORCES) {
            this.clearForces();
        }
        this.m_flags = ~World.LOCKED;
        // log.debug("ending step");
        this.m_profile.step.record(stepTimer.getMilliseconds());
    }

    clearForces() {
        for (let body = m_bodyList;
             body != null; body = body.getNext()) {
            body.m_force.setZero();
            body.m_torque = 0.0;
        }
    }

    drawDebugData() {
        if (this.m_debugDraw == null) {
            return;
        }
        let flags = m_debugDraw.getFlags();
        let wireframe = (flags & DebugDraw.e_wireframeDrawingBit) != 0;
        if ((flags & DebugDraw.e_shapeBit) != 0) {
            for (let b = m_bodyList;
                 b != null; b = b.getNext()) {
                xf.set(b.getTransform());
                for (let f = b.getFixtureList();
                     f != null; f = f.getNext()) {
                    if (b.isActive() == false) {
                        color.set(0.5, 0.5, 0.3);
                        this.drawShape(f, this.xf, this.color, wireframe);
                    } else if (b.getType() == BodyType.STATIC) {
                        color.set(0.5, 0.9, 0.3);
                        this.drawShape(f, this.xf, this.color, wireframe);
                    } else if (b.getType() == BodyType.KINEMATIC) {
                        color.set(0.5, 0.5, 0.9);
                        this.drawShape(f, this.xf, this.color, wireframe);
                    } else if (b.isAwake() == false) {
                        color.set(0.5, 0.5, 0.5);
                        this.drawShape(f, this.xf, this.color, wireframe);
                    } else {
                        color.set(0.9, 0.7, 0.7);
                        this.drawShape(f, this.xf, this.color, wireframe);
                    }
                }
            }
            this.drawParticleSystem(this.m_particleSystem);
        }
        if ((flags & DebugDraw.e_jointBit) != 0) {
            for (let j = m_jointList;
                 j != null; j = j.getNext()) {
                this.drawJoint(j);
            }
        }
        if ((flags & DebugDraw.e_pairBit) != 0) {
            color.set(0.3, 0.9, 0.9);
            for (let c = this.m_contactManager.m_contactList;
                 c != null; c = c.getNext()) {
                let fixtureA = c.getFixtureA();
                let fixtureB = c.getFixtureB();
                fixtureA.getAABB(c.getChildIndexA()).getCenterToOut(this.cA);
                fixtureB.getAABB(c.getChildIndexB()).getCenterToOut(this.cB);
                m_debugDraw.drawSegment(this.cA, this.cB, this.color);
            }
        }
        if ((flags & DebugDraw.e_aabbBit) != 0) {
            color.set(0.9, 0.3, 0.9);
            for (let b = m_bodyList;
                 b != null; b = b.getNext()) {
                if (b.isActive() == false) {
                    continue;
                }
                for (let f = b.getFixtureList();
                     f != null; f = f.getNext()) {
                    for (let i = 0;
                         i < f.m_proxyCount; ++i) {
                        let proxy = f.m_proxies[i];
                        let aabb = this.m_contactManager.m_broadPhase.getFatAABB(proxy.proxyId);
                        if (aabb != null) {
                            let vs = avs.get(4);
                            vs[0].set(aabb.lowerBound.x, aabb.lowerBound.y);
                            vs[1].set(aabb.upperBound.x, aabb.lowerBound.y);
                            vs[2].set(aabb.upperBound.x, aabb.upperBound.y);
                            vs[3].set(aabb.lowerBound.x, aabb.upperBound.y);
                            m_debugDraw.drawPolygon(vs, 4, this.color);
                        }
                    }
                }
            }
        }
        if ((flags & DebugDraw.e_centerOfMassBit) != 0) {
            for (let b = m_bodyList;
                 b != null; b = b.getNext()) {
                xf.set(b.getTransform());
                this.xf.p.set(b.getWorldCenter());
                m_debugDraw.drawTransform(this.xf);
            }
        }
        if ((flags & DebugDraw.e_dynamicTreeBit) != 0) {
            this.m_contactManager.m_broadPhase.drawTree(this.m_debugDraw);
        }
        m_debugDraw.flush();
    }

    queryAABB$QueryCallbackAABB(callback, aabb) {
        this.wqwrapper.broadPhase = this.m_contactManager.m_broadPhase;
        this.wqwrapper.callback = callback;
        this.m_contactManager.m_broadPhase.query(this.wqwrapper, aabb);
    }

    queryAABB$QueryCallbackParticleQueryCallbackAABB(callback, particleCallback, aabb) {
        this.wqwrapper.broadPhase = this.m_contactManager.m_broadPhase;
        this.wqwrapper.callback = callback;
        this.m_contactManager.m_broadPhase.query(this.wqwrapper, aabb);
        m_particleSystem.queryAABB(particleCallback, aabb);
    }

    queryAABB$ParticleQueryCallbackAABB(particleCallback, aabb) {
        m_particleSystem.queryAABB(particleCallback, aabb);
    }

    raycast$RayCastCallbackVec2Vec2(callback, point1, point2) {
        this.wrcwrapper.broadPhase = this.m_contactManager.m_broadPhase;
        this.wrcwrapper.callback = callback;
        this.input.maxFraction = 1.0;
        this.input.p1.set(point1);
        this.input.p2.set(point2);
        this.m_contactManager.m_broadPhase.raycast(this.wrcwrapper, this.input);
    }

    raycast$RayCastCallbackParticleRaycastCallbackVec2Vec2(callback, particleCallback, point1, point2) {
        this.wrcwrapper.broadPhase = this.m_contactManager.m_broadPhase;
        this.wrcwrapper.callback = callback;
        this.input.maxFraction = 1.0;
        this.input.p1.set(point1);
        this.input.p2.set(point2);
        this.m_contactManager.m_broadPhase.raycast(this.wrcwrapper, this.input);
        m_particleSystem.raycast(particleCallback, point1, point2);
    }

    raycast$ParticleRaycastCallbackVec2Vec2(particleCallback, point1, point2) {
        m_particleSystem.raycast(particleCallback, point1, point2);
    }

    getBodyList() {
        return this.m_bodyList;
    }

    getJointList() {
        return this.m_jointList;
    }

    getContactList() {
        return this.m_contactManager.m_contactList;
    }

    isSleepingAllowed() {
        return this.m_allowSleep;
    }

    setSleepingAllowed(sleepingAllowed) {
        this.m_allowSleep = sleepingAllowed;
    }

    setWarmStarting(flag) {
        this.m_warmStarting = flag;
    }

    isWarmStarting() {
        return this.m_warmStarting;
    }

    setContinuousPhysics(flag) {
        this.m_continuousPhysics = flag;
    }

    isContinuousPhysics() {
        return this.m_continuousPhysics;
    }

    getProxyCount() {
        return this.m_contactManager.m_broadPhase.getProxyCount();
    }

    getBodyCount() {
        return this.m_bodyCount;
    }

    getJointCount() {
        return this.m_jointCount;
    }

    getContactCount() {
        return this.m_contactManager.m_contactCount;
    }

    getTreeHeight() {
        return this.m_contactManager.m_broadPhase.getTreeHeight();
    }

    getTreeBalance() {
        return this.m_contactManager.m_broadPhase.getTreeBalance();
    }

    getTreeQuality() {
        return this.m_contactManager.m_broadPhase.getTreeQuality();
    }

    setGravity(gravity) {
        m_gravity.set(gravity);
    }

    getGravity() {
        return this.m_gravity;
    }

    isLocked() {
        return (this.m_flags & World.LOCKED) == World.LOCKED;
    }

    setAutoClearForces(flag) {
        if (flag) {
            this.m_flags = World.CLEAR_FORCES;
        } else {
            this.m_flags = ~World.CLEAR_FORCES;
        }
    }

    getAutoClearForces() {
        return (this.m_flags & World.CLEAR_FORCES) == World.CLEAR_FORCES;
    }

    getContactManager() {
        return this.m_contactManager;
    }

    getProfile() {
        return this.m_profile;
    }

    solve(step) {
        this.m_profile.solveInit.startAccum();
        this.m_profile.solveVelocity.startAccum();
        this.m_profile.solvePosition.startAccum();
        // update previous transforms
        for (let b = m_bodyList;
             b != null; b = b.m_next) {
            b.m_xf0.set(b.m_xf);
        }// Size the island for the worst case.
        island.init(this.m_bodyCount, this.m_contactManager.m_contactCount, this.m_jointCount, this.m_contactManager.m_contactListener);
        // Clear all the island flags.
        for (let b = m_bodyList;
             b != null; b = b.m_next) {
            b.m_flags = ~Body.e_islandFlag;
        }
        for (let c = this.m_contactManager.m_contactList;
             c != null; c = c.m_next) {
            c.m_flags = ~Contact.ISLAND_FLAG;
        }
        for (let j = m_jointList;
             j != null; j = j.m_next) {
            j.m_islandFlag = false;
        }// Build and simulate all awake islands.
        let stackSize = m_bodyCount;
        if (this.stack.length < stackSize) {
            this.stack = new Body[stackSize];
        }
        for (let seed = m_bodyList;
             seed != null; seed = seed.m_next) {
            if ((seed.m_flags & Body.e_islandFlag) == Body.e_islandFlag) {
                continue;
            }
            if (seed.isAwake() == false || seed.isActive() == false) {
                continue;
            }// The seed can be dynamic or kinematic.
            if (seed.getType() == BodyType.STATIC) {
                continue;
            }// Reset island and stack.
            island.clear();
            let stackCount = 0;
            this.stack[stackCount++] = seed;
            seed.m_flags = Body.e_islandFlag;
            // Perform a depth first search (DFS) on the constraint graph.
            while (stackCount > 0) {
                // Grab the next body off the stack and add it to the island.
                let b = this.stack[--stackCount];
                assert(b.isActive() == true);
                island.add(b);
                // Make sure the body is awake.
                b.setAwake(true);
                // To keep islands as small as possible, we don't
                // propagate islands across static bodies.
                if (b.getType() == BodyType.STATIC) {
                    continue;
                }// Search all contacts connected to this body.
                for (let ce = b.m_contactList;
                     ce != null; ce = ce.next) {
                    let contact = ce.contact;
                    // Has this contact already been added to an island?
                    if ((contact.m_flags & Contact.ISLAND_FLAG) == Contact.ISLAND_FLAG) {
                        continue;
                    }// Is this contact solid and touching?
                    if (contact.isEnabled() == false || contact.isTouching() == false) {
                        continue;
                    }// Skip sensors.
                    let sensorA = contact.m_fixtureA.m_isSensor;
                    let sensorB = contact.m_fixtureB.m_isSensor;
                    if (sensorA || sensorB) {
                        continue;
                    }
                    island.add(contact);
                    contact.m_flags = Contact.ISLAND_FLAG;
                    let other = ce.other;
                    // Was the other body already added to this island?
                    if ((other.m_flags & Body.e_islandFlag) == Body.e_islandFlag) {
                        continue;
                    }
                    assert(stackCount < stackSize);
                    this.stack[stackCount++] = other;
                    other.m_flags = Body.e_islandFlag;
                }// Search all joints connect to this body.
                for (let je = b.m_jointList;
                     je != null; je = je.next) {
                    if (je.joint.m_islandFlag == true) {
                        continue;
                    }
                    let other = je.other;
                    // Don't simulate joints connected to inactive bodies.
                    if (other.isActive() == false) {
                        continue;
                    }
                    island.add(je.joint);
                    je.joint.m_islandFlag = true;
                    if ((other.m_flags & Body.e_islandFlag) == Body.e_islandFlag) {
                        continue;
                    }
                    assert(stackCount < stackSize);
                    this.stack[stackCount++] = other;
                    other.m_flags = Body.e_islandFlag;
                }
            }
            island.solve(this.m_profile, step, this.m_gravity, this.m_allowSleep);
            // Post solve cleanup.
            for (let i = 0;
                 i < this.island.m_bodyCount; ++i) {
                // Allow static bodies to participate in other islands.
                let b = this.island.m_bodies[i];
                if (b.getType() == BodyType.STATIC) {
                    b.m_flags = ~Body.e_islandFlag;
                }
            }
        }
        this.m_profile.solveInit.endAccum();
        this.m_profile.solveVelocity.endAccum();
        this.m_profile.solvePosition.endAccum();
        broadphaseTimer.reset();
        // Synchronize fixtures, check for out of range bodies.
        for (let b = m_bodyList;
             b != null; b = b.getNext()) {
            // If a body was not in an island then it did not move.
            if ((b.m_flags & Body.e_islandFlag) == 0) {
                continue;
            }
            if (b.getType() == BodyType.STATIC) {
                continue;
            }// Update fixtures (for broad-phase).
            b.synchronizeFixtures();
        }// Look for new contacts.
        m_contactManager.findNewContacts();
        this.m_profile.broadphase.record(broadphaseTimer.getMilliseconds());
    }

    solveTOI(step) {
        const island = toiIsland;
        island.init(2 * Settings.maxTOIContacts, Settings.maxTOIContacts, 0, this.m_contactManager.m_contactListener);
        if (this.m_stepComplete) {
            for (let b = m_bodyList;
                 b != null; b = b.m_next) {
                b.m_flags = ~Body.e_islandFlag;
                b.m_sweep.alpha0 = 0.0;
            }
            for (let c = this.m_contactManager.m_contactList;
                 c != null; c = c.m_next) {
                // Invalidate TOI
                c.m_flags = ~(Contact.TOI_FLAG | Contact.ISLAND_FLAG);
                c.m_toiCount = 0;
                c.m_toi = 1.0;
            }
        }// Find TOI events and solve them.
        for (; ;) {
            // Find the first TOI.
            let minContact = null;
            let minAlpha = 1.0;
            for (let c = this.m_contactManager.m_contactList;
                 c != null; c = c.m_next) {
                // Is this contact disabled?
                if (c.isEnabled() == false) {
                    continue;
                }// Prevent excessive sub-stepping.
                if (c.m_toiCount > Settings.maxSubSteps) {
                    continue;
                }
                let alpha = 1.0;
                if ((c.m_flags & Contact.TOI_FLAG) != 0) {
                    // This contact has a valid cached TOI.
                    alpha = c.m_toi;
                } else {
                    let fA = c.getFixtureA();
                    let fB = c.getFixtureB();
                    // Is there a sensor?
                    if (fA.isSensor() || fB.isSensor()) {
                        continue;
                    }
                    let bA = fA.getBody();
                    let bB = fB.getBody();
                    let typeA = bA.m_type;
                    let typeB = bB.m_type;
                    assert(typeA == BodyType.DYNAMIC || typeB == BodyType.DYNAMIC);
                    let activeA = bA.isAwake() && typeA != BodyType.STATIC;
                    let activeB = bB.isAwake() && typeB != BodyType.STATIC;
                    // Is at least one body active (awake and dynamic or kinematic)?
                    if (activeA == false && activeB == false) {
                        continue;
                    }
                    let collideA = bA.isBullet() || typeA != BodyType.DYNAMIC;
                    let collideB = bB.isBullet() || typeB != BodyType.DYNAMIC;
                    // Are these two non-bullet dynamic bodies?
                    if (collideA == false && collideB == false) {
                        continue;
                    }// Compute the TOI for this contact.
                    // Put the sweeps onto the same time interval.
                    let alpha0 = bA.m_sweep.alpha0;
                    if (bA.m_sweep.alpha0 < bB.m_sweep.alpha0) {
                        alpha0 = bB.m_sweep.alpha0;
                        bA.m_sweep.advance(alpha0);
                    } else if (bB.m_sweep.alpha0 < bA.m_sweep.alpha0) {
                        alpha0 = bA.m_sweep.alpha0;
                        bB.m_sweep.advance(alpha0);
                    }
                    assert(alpha0 < 1.0);
                    let indexA = c.getChildIndexA();
                    let indexB = c.getChildIndexB();
                    // Compute the time of impact in interval [0, minTOI]
                    const input = toiInput;
                    input.proxyA.set(fA.getShape(), indexA);
                    input.proxyB.set(fB.getShape(), indexB);
                    input.sweepA.set(bA.m_sweep);
                    input.sweepB.set(bB.m_sweep);
                    input.tMax = 1.0;
                    pool.getTimeOfImpact().timeOfImpact(this.toiOutput, input);
                    // Beta is the fraction of the remaining portion of the .
                    let beta = this.toiOutput.t;
                    if (this.toiOutput.state == TOIOutputState.TOUCHING) {
                        alpha = MathUtils.min(alpha0 + (1.0 - alpha0) * beta, 1.0);
                    } else {
                        alpha = 1.0;
                    }
                    c.m_toi = alpha;
                    c.m_flags = Contact.TOI_FLAG;
                }
                if (alpha < minAlpha) {
                    // This is the minimum TOI found so far.
                    minContact = c;
                    minAlpha = alpha;
                }
            }
            if (minContact == null || 1.0 - 10.0 * Settings.EPSILON < minAlpha) {
                // No more TOI events. Done!
                this.m_stepComplete = true;
                break;
            }// Advance the bodies to the TOI.
            let fA = minContact.getFixtureA();
            let fB = minContact.getFixtureB();
            let bA = fA.getBody();
            let bB = fB.getBody();
            backup1.set(bA.m_sweep);
            backup2.set(bB.m_sweep);
            bA.advance(minAlpha);
            bB.advance(minAlpha);
            // The TOI contact likely has some new contact points.
            minContact.update(this.m_contactManager.m_contactListener);
            minContact.m_flags = ~Contact.TOI_FLAG;
            ++minContact.m_toiCount;
            // Is the contact solid?
            if (minContact.isEnabled() == false || minContact.isTouching() == false) {
                // Restore the sweeps.
                minContact.setEnabled(false);
                bA.m_sweep.set(this.backup1);
                bB.m_sweep.set(this.backup2);
                bA.synchronizeTransform();
                bB.synchronizeTransform();
                continue;
            }
            bA.setAwake(true);
            bB.setAwake(true);
            // Build the island
            island.clear();
            island.add(bA);
            island.add(bB);
            island.add(minContact);
            bA.m_flags = Body.e_islandFlag;
            bB.m_flags = Body.e_islandFlag;
            minContact.m_flags = Contact.ISLAND_FLAG;
            // Get contacts on bodyA and bodyB.
            this.tempBodies[0] = bA;
            this.tempBodies[1] = bB;
            for (let i = 0;
                 i < 2; ++i) {
                let body = this.tempBodies[i];
                if (body.m_type == BodyType.DYNAMIC) {
                    for (let ce = body.m_contactList;
                         ce != null; ce = ce.next) {
                        if (island.m_bodyCount == island.m_bodyCapacity) {
                            break;
                        }
                        if (island.m_contactCount == island.m_contactCapacity) {
                            break;
                        }
                        let contact = ce.contact;
                        // Has this contact already been added to the island?
                        if ((contact.m_flags & Contact.ISLAND_FLAG) != 0) {
                            continue;
                        }// Only add static, kinematic, or bullet bodies.
                        let other = ce.other;
                        if (other.m_type == BodyType.DYNAMIC && body.isBullet() == false && other.isBullet() == false) {
                            continue;
                        }// Skip sensors.
                        let sensorA = contact.m_fixtureA.m_isSensor;
                        let sensorB = contact.m_fixtureB.m_isSensor;
                        if (sensorA || sensorB) {
                            continue;
                        }// Tentatively advance the body to the TOI.
                        backup1.set(other.m_sweep);
                        if ((other.m_flags & Body.e_islandFlag) == 0) {
                            other.advance(minAlpha);
                        }// Update the contact points
                        contact.update(this.m_contactManager.m_contactListener);
                        // Was the contact disabled by the user?
                        if (contact.isEnabled() == false) {
                            other.m_sweep.set(this.backup1);
                            other.synchronizeTransform();
                            continue;
                        }// Are there contact points?
                        if (contact.isTouching() == false) {
                            other.m_sweep.set(this.backup1);
                            other.synchronizeTransform();
                            continue;
                        }// Add the contact to the island
                        contact.m_flags = Contact.ISLAND_FLAG;
                        island.add(contact);
                        // Has the other body already been added to the island?
                        if ((other.m_flags & Body.e_islandFlag) != 0) {
                            continue;
                        }// Add the other body to the island.
                        other.m_flags = Body.e_islandFlag;
                        if (other.m_type != BodyType.STATIC) {
                            other.setAwake(true);
                        }
                        island.add(other);
                    }
                }
            }
            this.subStep.dt = (1.0 - minAlpha) * step.dt;
            this.subStep.inv_dt = 1.0 / this.subStep.dt;
            this.subStep.dtRatio = 1.0;
            this.subStep.positionIterations = 20;
            this.subStep.velocityIterations = step.velocityIterations;
            this.subStep.warmStarting = false;
            island.solveTOI(this.subStep, bA.m_islandIndex, bB.m_islandIndex);
            // Reset island flags and synchronize broad-phase proxies.
            for (let i = 0;
                 i < island.m_bodyCount; ++i) {
                let body = island.m_bodies[i];
                body.m_flags = ~Body.e_islandFlag;
                if (body.m_type != BodyType.DYNAMIC) {
                    continue;
                }
                body.synchronizeFixtures();
                // Invalidate all contact TOIs on this displaced body.
                for (let ce = body.m_contactList;
                     ce != null; ce = ce.next) {
                    ce.contact.m_flags = ~(Contact.TOI_FLAG | Contact.ISLAND_FLAG);
                }
            }// Commit fixture proxy movements to the broad-phase so that new contacts are created.
            // Also, some contacts can be destroyed.
            m_contactManager.findNewContacts();
            if (this.m_subStepping) {
                this.m_stepComplete = false;
                break;
            }
        }
    }

    drawJoint(joint) {
        let bodyA = joint.getBodyA();
        let bodyB = joint.getBodyB();
        let xf1 = bodyA.getTransform();
        let xf2 = bodyB.getTransform();
        let x1 = xf1.p;
        let x2 = xf2.p;
        let p1 = pool.popVec2();
        let p2 = pool.popVec2();
        joint.getAnchorA(p1);
        joint.getAnchorB(p2);
        color.set(0.5, 0.8, 0.8);
        switch (joint.getType()) {
            // TODO djm write after writing joints
            case DISTANCE:
                m_debugDraw.drawSegment(p1, p2, this.color);
                break;
            case PULLEY: {
                let pulley = joint;
                let s1 = pulley.getGroundAnchorA();
                let s2 = pulley.getGroundAnchorB();
                m_debugDraw.drawSegment(s1, p1, this.color);
                m_debugDraw.drawSegment(s2, p2, this.color);
                m_debugDraw.drawSegment(s1, s2, this.color);
            }
                break;
            case CONSTANT_VOLUME:
            case MOUSE:// don't draw this
                break;
            default:
                m_debugDraw.drawSegment(x1, p1, this.color);
                m_debugDraw.drawSegment(p1, p2, this.color);
                m_debugDraw.drawSegment(x2, p2, this.color);
        }
        pool.pushVec2(2);
    }

    drawShape(fixture, xf, color, wireframe) {
        switch (fixture.getType()) {
            case CIRCLE: {
                let circle = fixture.getShape();
                // Vec2 center = Mul(xf, circle.m_p);
                Transform.mulToOutUnsafe(xf, circle.m_p, this.center);
                let radius = circle.m_radius;
                xf.q.getXAxis(this.axis);
                if (fixture.getUserData() != null && fixture.getUserData().equals(World.LIQUID_INT)) {
                    let b = fixture.getBody();
                    liquidOffset.set(b.m_linearVelocity);
                    let linVelLength = b.m_linearVelocity.length();
                    if (this.averageLinearVel == -1) {
                        this.averageLinearVel = linVelLength;
                    } else {
                        this.averageLinearVel = .98 * this.averageLinearVel + .02 * linVelLength;
                    }
                    liquidOffset.mulLocal(this.liquidLength / this.averageLinearVel / 2);
                    circCenterMoved.set(this.center).addLocal(this.liquidOffset);
                    center.subLocal(this.liquidOffset);
                    m_debugDraw.drawSegment(this.center, this.circCenterMoved, this.liquidColor);
                    return;
                }
                if (wireframe) {
                    m_debugDraw.drawCircle(this.center, radius, this.axis, color);
                } else {
                    m_debugDraw.drawSolidCircle(this.center, radius, this.axis, color);
                }
            }
                break;
            case POLYGON: {
                let poly = fixture.getShape();
                let vertexCount = poly.m_count;
                assert(vertexCount <= Settings.maxPolygonVertices);
                let vertices = tlvertices.get(Settings.maxPolygonVertices);
                for (let i = 0;
                     i < vertexCount; ++i) {
                    // vertices[i] = Mul(xf, poly.m_vertices[i]);
                    Transform.mulToOutUnsafe(xf, poly.m_vertices[i], vertices[i]);
                }
                if (wireframe) {
                    m_debugDraw.drawPolygon(vertices, vertexCount, color);
                } else {
                    m_debugDraw.drawSolidPolygon(vertices, vertexCount, color);
                }
            }
                break;
            case EDGE: {
                let edge = fixture.getShape();
                Transform.mulToOutUnsafe(xf, edge.m_vertex1, this.v1);
                Transform.mulToOutUnsafe(xf, edge.m_vertex2, this.v2);
                m_debugDraw.drawSegment(this.v1, this.v2, color);
            }
                break;
            case CHAIN: {
                let chain = fixture.getShape();
                let count = chain.m_count;
                let vertices = chain.m_vertices;
                Transform.mulToOutUnsafe(xf, vertices[0], this.v1);
                for (let i = 1;
                     i < count; ++i) {
                    Transform.mulToOutUnsafe(xf, vertices[i], this.v2);
                    m_debugDraw.drawSegment(this.v1, this.v2, color);
                    m_debugDraw.drawCircle(this.v1, 0.05, color);
                    v1.set(this.v2);
                }
            }
                break;
            default:
                break;
        }
    }

    drawParticleSystem(system) {
        let wireframe = (m_debugDraw.getFlags() & DebugDraw.e_wireframeDrawingBit) != 0;
        let particleCount = system.getParticleCount();
        if (particleCount != 0) {
            let particleRadius = system.getParticleRadius();
            let positionBuffer = system.getParticlePositionBuffer();
            let colorBuffer = null;
            if (system.m_colorBuffer.data != null) {
                colorBuffer = system.getParticleColorBuffer();
            }
            if (wireframe) {
                m_debugDraw.drawParticlesWireframe(positionBuffer, particleRadius, colorBuffer, particleCount);
            } else {
                m_debugDraw.drawParticles(positionBuffer, particleRadius, colorBuffer, particleCount);
            }
        }
    }

    createParticle(def) {
        assert(this.isLocked() == false);
        if (this.isLocked()) {
            return 0;
        }
        let p = m_particleSystem.createParticle(def);
        return p;
    }

    destroyParticle$int(index) {
        this.destroyParticle(index, false);
    }

    destroyParticle$intboolean(index, callDestructionListener) {
        m_particleSystem.destroyParticle(index, callDestructionListener);
    }

    destroyParticlesInShape$ShapeTransform(shape, xf) {
        return this.destroyParticlesInShape(shape, xf, false);
    }

    destroyParticlesInShape$ShapeTransformboolean(shape, xf, callDestructionListener) {
        assert(this.isLocked() == false);
        if (this.isLocked()) {
            return 0;
        }
        return m_particleSystem.destroyParticlesInShape(shape, xf, callDestructionListener);
    }

    createParticleGroup(def) {
        assert(this.isLocked() == false);
        if (this.isLocked()) {
            return null;
        }
        let g = m_particleSystem.createParticleGroup(def);
        return g;
    }

    joinParticleGroups(groupA, groupB) {
        assert(this.isLocked() == false);
        if (this.isLocked()) {
            return;
        }
        m_particleSystem.joinParticleGroups(groupA, groupB);
    }

    destroyParticlesInGroup$ParticleGroupboolean(group, callDestructionListener) {
        assert(this.isLocked() == false);
        if (this.isLocked()) {
            return;
        }
        m_particleSystem.destroyParticlesInGroup(group, callDestructionListener);
    }

    destroyParticlesInGroup$ParticleGroup(group) {
        this.destroyParticlesInGroup(group, false);
    }

    getParticleGroupList() {
        return m_particleSystem.getParticleGroupList();
    }

    getParticleGroupCount() {
        return m_particleSystem.getParticleGroupCount();
    }

    getParticleCount() {
        return m_particleSystem.getParticleCount();
    }

    getParticleMaxCount() {
        return m_particleSystem.getParticleMaxCount();
    }

    setParticleMaxCount(count) {
        m_particleSystem.setParticleMaxCount(count);
    }

    setParticleDensity(density) {
        m_particleSystem.setParticleDensity(density);
    }

    getParticleDensity() {
        return m_particleSystem.getParticleDensity();
    }

    setParticleGravityScale(gravityScale) {
        m_particleSystem.setParticleGravityScale(gravityScale);
    }

    getParticleGravityScale() {
        return m_particleSystem.getParticleGravityScale();
    }

    setParticleDamping(damping) {
        m_particleSystem.setParticleDamping(damping);
    }

    getParticleDamping() {
        return m_particleSystem.getParticleDamping();
    }

    setParticleRadius(radius) {
        m_particleSystem.setParticleRadius(radius);
    }

    getParticleRadius() {
        return m_particleSystem.getParticleRadius();
    }

    getParticleFlagsBuffer() {
        return m_particleSystem.getParticleFlagsBuffer();
    }

    getParticlePositionBuffer() {
        return m_particleSystem.getParticlePositionBuffer();
    }

    getParticleVelocityBuffer() {
        return m_particleSystem.getParticleVelocityBuffer();
    }

    getParticleColorBuffer() {
        return m_particleSystem.getParticleColorBuffer();
    }

    getParticleGroupBuffer() {
        return m_particleSystem.getParticleGroupBuffer();
    }

    getParticleUserDataBuffer() {
        return m_particleSystem.getParticleUserDataBuffer();
    }

    setParticleFlagsBuffer(buffer, capacity) {
        m_particleSystem.setParticleFlagsBuffer(buffer, capacity);
    }

    setParticlePositionBuffer(buffer, capacity) {
        m_particleSystem.setParticlePositionBuffer(buffer, capacity);
    }

    setParticleVelocityBuffer(buffer, capacity) {
        m_particleSystem.setParticleVelocityBuffer(buffer, capacity);
    }

    setParticleColorBuffer(buffer, capacity) {
        m_particleSystem.setParticleColorBuffer(buffer, capacity);
    }

    setParticleUserDataBuffer(buffer, capacity) {
        m_particleSystem.setParticleUserDataBuffer(buffer, capacity);
    }

    getParticleContacts() {
        return this.m_particleSystem.m_contactBuffer;
    }

    getParticleContactCount() {
        return this.m_particleSystem.m_contactCount;
    }

    getParticleBodyContacts() {
        return this.m_particleSystem.m_bodyContactBuffer;
    }

    getParticleBodyContactCount() {
        return this.m_particleSystem.m_bodyContactCount;
    }

    computeParticleCollisionEnergy() {
        return m_particleSystem.computeParticleCollisionEnergy();
    }
}

export class WorldQueryWrapper {

    /**
     * @var {BroadPhase}
     */
    broadPhase;

    /**
     * @var {QueryCallback}
     */
    callback;


    treeCallback(nodeId) {
        let proxy = broadPhase.getUserData(nodeId);
        return callback.reportFixture(proxy.fixture);
    }
}

export class WorldRayCastWrapper {

    /**
     * @var {RayCastOutput}
     */
    output = new RayCastOutput();

    /**
     * @var {Vec2}
     */
    temp = new Vec2();

    /**
     * @var {Vec2}
     */
    point = new Vec2();

    /**
     * @var {BroadPhase}
     */
    broadPhase;

    /**
     * @var {RayCastCallback}
     */
    callback;


    raycastCallback(input, nodeId) {
        let userData = broadPhase.getUserData(nodeId);
        let proxy = userData;
        let fixture = proxy.fixture;
        let index = proxy.childIndex;
        let hit = fixture.raycast(this.output, input, index);
        if (hit) {
            let fraction = this.output.fraction;
            // Vec2 point = (1.0f - fraction) * input.p1 + fraction * input.p2;
            temp.set(input.p2).mulLocal(fraction);
            point.set(input.p1).mulLocal(1 - fraction).addLocal(this.temp);
            return callback.reportFixture(fixture, this.point, this.output.normal, fraction);
        }
        return input.maxFraction;
    }
}
