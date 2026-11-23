import { CheckCircle2, Info } from '@Components/icons';
import { getForgotPasswordNoticeAppearance } from '../forgotPassword.noticeAppearance';
describe('getForgotPasswordNoticeAppearance', () => {
  it('provides the appropriate icon for success and information', () => {
    expect(getForgotPasswordNoticeAppearance(true).NoticeIcon).toBe(CheckCircle2);
    expect(getForgotPasswordNoticeAppearance(false).NoticeIcon).toBe(Info);
  });
});
