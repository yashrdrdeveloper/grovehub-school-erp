import { useEffect, useRef } from 'react';

// WebGL2 Soffit gradient. One fullscreen triangle; no third-party libraries.
// The palette and animation settings match the supplied Soffit specification.
const CONFIG = {
  bgColor: '#062334', colorA: '#2b7be8', colorB: '#2bc8e8', colorC: '#4fe0d8', colorD: '#7fe8a0',
  ambient: .24, amount: .2, bounce: .38, bounceCurve: 4.25, breathe: .29, contrast: 2.45,
  curve: 3.32, direct: .97, dither: .58, flow: .475, glow: .38, horizon: .36,
  lacunarity: 1.99, lift: .11, midpoint: .57, moteScale: 7, motes: .074,
  parallax: .0137, rock: .12, roughness: .29, scale: 1, sink: .24, speed: .33,
  spillCentre: .3, spillFloor: .26, spillWidth: 2.18, spread: .41,
  steer: -.13, sweep: .5, tilt: 1.87, vignette: .21, warp: 2.58, warpScale: .78,
};
const VERT = `#version 300 es
void main(){vec2 p=vec2((gl_VertexID<<1)&2,gl_VertexID&2);gl_Position=vec4(p*2.0-1.0,0.0,1.0);}`;
const FRAG = `#version 300 es
precision highp float;
out vec4 fragColor;
uniform vec2 iResolution, iMouse;
uniform float iTime, uScale;
uniform vec3 uBg,uColorA,uColorB,uColorC,uColorD;
uniform float uSpeed,uTilt,uRock,uHorizon,uBreathe,uSpread,uCurve,uDirect,uBounce,uBounceCurve;
uniform float uSpillCentre,uSpillWidth,uSpillFloor,uAmount,uWarp,uWarpScale,uFlow,uRoughness,uLacunarity,uMotes,uMoteScale;
uniform float uAmbient,uContrast,uMidpoint,uSink,uGlow,uDither,uVignette,uSteer,uLift,uSweep,uParallax;
vec2 hash2(vec2 p){p=vec2(dot(p,vec2(127.1,311.7)),dot(p,vec2(269.5,183.3)));return -1.0+2.0*fract(sin(p)*43758.5453123);}
float snoise(vec2 p){const float K1=.366025404,K2=.211324865;vec2 i=floor(p+(p.x+p.y)*K1);vec2 a=p-i+(i.x+i.y)*K2;float m=step(a.y,a.x);vec2 o=vec2(m,1.0-m);vec2 b=a-o+K2;vec2 c=a-1.0+2.0*K2;vec3 h=max(.5-vec3(dot(a,a),dot(b,b),dot(c,c)),0.0);vec3 n=h*h*h*h*vec3(dot(a,hash2(i)),dot(b,hash2(i+o)),dot(c,hash2(i+1.0)));return dot(n,vec3(70.0));}
float fbm(vec2 p){float v=0.0,amp=.5;for(int i=0;i<4;i++){v+=amp*snoise(p);p*=uLacunarity;amp*=uRoughness;}return v;}
vec3 ramp4(float t){vec3 c=mix(uColorA,uColorB,smoothstep(0.0,.36,t));c=mix(c,uColorC,smoothstep(.32,.70,t));c=mix(c,uColorD,smoothstep(.66,1.0,t));return c;}
float triDither(vec2 fc){float a=fract(sin(dot(fc,vec2(12.9898,78.233)))*43758.5453);float b=fract(sin(dot(fc+17.0,vec2(12.9898,78.233)))*43758.5453);return(a+b-1.0)/255.0;}
void main(){
 vec2 uv=(gl_FragCoord.xy-.5*iResolution)/iResolution.y;uv*=uScale;
 vec2 iM=iMouse*uScale;float t=iTime*uSpeed;vec2 p=uv-iM*uParallax;
 float tilt=uTilt+sin(t*.13)*uRock+iM.x*uSteer;vec2 dir=vec2(cos(tilt),sin(tilt));
 float axis=dot(p,dir),across=dot(p,vec2(-dir.y,dir.x));
 vec2 q=vec2(fbm(p*uWarpScale+vec2(0.0,t*uFlow)),fbm(p*uWarpScale+vec2(5.2,1.3)-t*uFlow*.7));
 float air=fbm(p+uWarp*q+vec2(t*.12,-t*.09))*.5+.5;
 float horizon=uHorizon+sin(t*.09+2.1)*uBreathe-iM.y*uLift;
 float alt=clamp(.5+(axis-horizon)*uSpread+(air-.5)*uAmount,0.0,1.0);
 float ac=(across-uSpillCentre-iM.x*uSweep)/max(.05,uSpillWidth);
 float spill=mix(uSpillFloor,1.0,exp(-ac*ac));
 float direct=pow(alt,max(.05,uCurve))*uDirect*spill;
 float bounce=uBounce*pow(1.0-alt,max(.05,uBounceCurve));
 float f=uAmbient+direct+bounce;
 f+=uMotes*snoise(p*uMoteScale+vec2(-t*.5,t*.35))*.5*alt;
 f=clamp((f-uMidpoint)*uContrast+.5,0.0,1.0);
 vec3 col=ramp4(f);col+=uColorD*uGlow*pow(f,4.0);
 col=mix(uBg,col,smoothstep(0.0,max(.01,uSink),f)*.90+.10);
 col*=1.0-uVignette*dot(uv,uv);col+=triDither(gl_FragCoord.xy)*uDither;
 fragColor=vec4(clamp(col,0.0,1.0),1.0);
}`;

export default function SoffitGradient() {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const gl = canvas.getContext('webgl2', { alpha: false, antialias: false, depth: false, stencil: false, powerPreference: 'high-performance' });
    if (!gl) return; // The CSS background is the fallback.
    const shader = (type: number, source: string) => {
      const s = gl.createShader(type);
      if (!s) throw new Error('Cannot create WebGL shader');
      gl.shaderSource(s, source); gl.compileShader(s);
      if (!gl.getShaderParameter(s, gl.COMPILE_STATUS)) {
        const msg = gl.getShaderInfoLog(s); gl.deleteShader(s); throw new Error(msg || 'Shader compilation failed');
      }
      return s;
    };
    let program: WebGLProgram | null = null;
    let vertex: WebGLShader | null = null;
    let fragment: WebGLShader | null = null;
    let vao: WebGLVertexArrayObject | null = null;
    try {
      vertex = shader(gl.VERTEX_SHADER, VERT);
      fragment = shader(gl.FRAGMENT_SHADER, FRAG);
      program = gl.createProgram();
      if (!program) throw new Error('Cannot create WebGL program');
      gl.attachShader(program, vertex); gl.attachShader(program, fragment); gl.linkProgram(program);
      if (!gl.getProgramParameter(program, gl.LINK_STATUS)) throw new Error(gl.getProgramInfoLog(program) || 'Shader link failed');
      gl.useProgram(program);
      vao = gl.createVertexArray();
      gl.bindVertexArray(vao);
    } catch (err) {
      console.warn('Soffit WebGL unavailable; using fallback', err);
      if (vao) gl.deleteVertexArray(vao);
      if (program) gl.deleteProgram(program);
      if (vertex) gl.deleteShader(vertex);
      if (fragment) gl.deleteShader(fragment);
      return;
    }
    const locations = new Map<string, WebGLUniformLocation | null>();
    const loc = (name: string) => {
      if (!locations.has(name)) locations.set(name, gl.getUniformLocation(program!, name));
      return locations.get(name)!;
    };
    const one = (n: string, v: number) => gl.uniform1f(loc(n), v);
    const two = (n: string, x: number, y: number) => gl.uniform2f(loc(n), x, y);
    const rgb = (n: string, hex: string) => {
      const raw = hex.slice(1);
      gl.uniform3f(loc(n), parseInt(raw.slice(0, 2), 16)/255, parseInt(raw.slice(2, 4), 16)/255, parseInt(raw.slice(4, 6), 16)/255);
    };
    rgb('uBg', CONFIG.bgColor);
    for (const key of ['colorA','colorB','colorC','colorD'] as const) rgb('u'+key.charAt(0).toUpperCase()+key.slice(1), CONFIG[key]);
    for (const [key, value] of Object.entries(CONFIG)) if (typeof value === 'number') one('u'+key.charAt(0).toUpperCase()+key.slice(1), value);
    let width = 1, height = 1;
    const resize = () => {
      const r = canvas.getBoundingClientRect();
      width = Math.max(1, Math.round(r.width)); height = Math.max(1, Math.round(r.height));
      if (canvas.width !== width || canvas.height !== height) { canvas.width = width; canvas.height = height; }
      gl.viewport(0,0,width,height); two('iResolution', width, height);
    };
    resize();
    const observer = typeof ResizeObserver !== 'undefined' ? new ResizeObserver(resize) : null;
    observer?.observe(canvas);
    const motion = window.matchMedia?.('(prefers-reduced-motion: reduce)');
    const mouse = {x:0,y:0,ax:0,ay:0,tx:0,ty:0};
    const aim = (event: PointerEvent) => {
      const r = canvas.getBoundingClientRect();
      mouse.tx = ((event.clientX - r.left)/Math.max(r.width,1)-.5)*(r.width/Math.max(r.height,1));
      mouse.ty = .5-(event.clientY-r.top)/Math.max(r.height,1);
    };
    canvas.addEventListener('pointermove',aim,{passive:true});
    canvas.addEventListener('pointerdown',aim,{passive:true});
    let raf = 0, clock = 0, last = 0;
    const frame = (now: number) => {
      raf = requestAnimationFrame(frame);
      if (!last) last=now;
      const dt=Math.min(50,Math.max(4.167,now-last)); last=now;
      if (document.hidden || !canvas.isConnected) return;
      if (!motion?.matches) {
        clock+=dt*.001;
        const step=Math.min(2.2,dt*.06);
        mouse.ax+=(mouse.tx-mouse.ax)*.105*step;mouse.ay+=(mouse.ty-mouse.ay)*.105*step;
        mouse.x+=(mouse.ax-mouse.x)*.043*step;mouse.y+=(mouse.ay-mouse.y)*.043*step;
      }
      one('iTime',clock);two('iMouse',mouse.x,mouse.y);
      gl.drawArrays(gl.TRIANGLES,0,3);
    };
    frame(performance.now());
    return () => {
      cancelAnimationFrame(raf);observer?.disconnect();
      canvas.removeEventListener('pointermove',aim);canvas.removeEventListener('pointerdown',aim);
      gl.bindVertexArray(null);if (vao) gl.deleteVertexArray(vao);
      if (program) gl.deleteProgram(program);
      if (vertex) gl.deleteShader(vertex);if (fragment) gl.deleteShader(fragment);
    };
  }, []);
  return <canvas ref={canvasRef} className="lv2-gradient" aria-hidden="true" style={{ pointerEvents: 'auto' }} />;
}
