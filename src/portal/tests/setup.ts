import Enzyme from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';

Enzyme.configure({ adapter: new Adapter() });
const jsdom = require('jsdom');
const { JSDOM } = jsdom;
const { window } = new JSDOM(`<!doctype html><html><head></head><body></body></html>`);
const { document } = (new JSDOM(``)).window;

global.document = document;
global.window = window;
global.navigator = window.navigator;
global.XMLHttpRequest = window.XMLHttpRequest;
global.getComputedStyle = window.getComputedStyle;
global.HTMLElement = window.HTMLElement;

Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: jest.fn().mockImplementation((query) => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: jest.fn(),
        removeListener: jest.fn(),
        addEventListener: jest.fn(),
        removeEventListener: jest.fn(),
        dispatchEvent: jest.fn(),
    })),
});

function mockRouter() {
    const original = jest.requireActual('react-router-dom');
    return {
      ...original,
      useHistory: jest.fn().mockReturnValue({
        location: {
            pathname: '/home',
            key: "k0pble",
            hash: '',
            search: '',
            state: null,
        },
        push: function (path: string){
          this.location.pathname = path;
        },
      }),
    };
}

function mockUmi() {
    const original = jest.requireActual('umi');
    return {
      ...original,
      useSelector: (callback: any) => {
          return callback;
      },
      ComputeUnitRadioAction: {
        GET_DATA: 'computeUnitRadio/getComputeUnit',
      },
      CommonActions: {
        INIT_USER: 'common/initUser',
        UPDATE_STATE: 'common/updateState',
      },
      MessageDetailAction: {
        INIT_DATA: 'messageDetail/initData',
        SIGN_UP: 'messageDetail/signUp',
        INIT_ENV: 'messageDetail/initEnv',
        UPDATE_STATE: 'messageDetail/updateState'
      }
    };
}

jest.mock('react-router-dom', () => mockRouter());
jest.mock('umi', () => mockUmi());